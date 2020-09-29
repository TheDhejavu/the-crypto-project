package p2p

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"sync"
	"syscall"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
	discovery "github.com/libp2p/go-libp2p-discovery"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	mplex "github.com/libp2p/go-libp2p-mplex"
	yamux "github.com/libp2p/go-libp2p-yamux"
	tcp "github.com/libp2p/go-tcp-transport"
	ws "github.com/libp2p/go-ws-transport"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"gopkg.in/vrecan/death.v3"

	"github.com/libp2p/go-libp2p-core/host"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	blockchain "github.com/workspace/the-crypto-project/core"
	"github.com/workspace/the-crypto-project/memopool"
)

type Network struct {
	Host           host.Host
	GeneralChannel *Channel
	MiningChannel  *Channel
	Blockchain     *blockchain.Blockchain
}

type Version struct {
	Version    int
	BestHeight int
	SendFrom   string
}

type GetBlocks struct {
	SendFrom string
	Height   int
}

type Block struct {
	SendFrom string
	Block    []byte
}

type GetData struct {
	SendFrom string
	Type     string
	ID       []byte
}

type Inv struct {
	SendFrom string
	Type     string
	Items    [][]byte
}

const (
	version       = 1
	commandLength = 12
)

var (
	GeneralChannel  = "general-channel"
	MiningChannel   = "mining-channel"
	MinerAddress    = ""
	blocksInTransit = [][]byte{}
	memoryPool      = new(memopool.MemoPool)
)

func CmdToBytes(cmd string) []byte {
	var bytes [commandLength]byte
	for i, c := range cmd {
		bytes[i] = byte(c)
	}
	return bytes[:]
}

func BytesToCmd(bytes []byte) string {
	var cmd []byte
	for _, b := range bytes {
		if b != byte(0) {
			cmd = append(cmd, b)
		}
	}
	return fmt.Sprintf("%s", cmd)
}

func GobEncode(data interface{}) []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(data)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

func CloseDB(chain *blockchain.Blockchain) {
	d := death.NewDeath(syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	d.WaitForDeathWithFunc(func() {
		defer os.Exit(1)
		defer runtime.Goexit()
		chain.Database.Close()
	})
}

func (net *Network) SendBlock(peerId string, b *blockchain.Block) {
	data := Block{net.Host.ID().Pretty(), b.Serialize()}
	payload := GobEncode(data)
	request := append(CmdToBytes("block"), payload...)
	if peerId != "" {
		net.GeneralChannel.Publish("Recieved block", request, peerId)
	} else {
		net.GeneralChannel.Publish("Recieved block", request, "")
	}
}

func (net *Network) HandleBlocks(content *ChannelContent) {
	var buff bytes.Buffer
	var payload Block

	buff.Write(content.Payload[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)

	if err != nil {
		log.Panic(err)
	}

	blockData := payload.Block
	block := blockchain.DeSerialize(blockData)

	// fmt.Printf("Valid: %s\n", strconv.FormatBool(validate))
	// Verify block before adding it to the blockchain
	if block.IsGenesis() {
		net.Blockchain.AddBlock(block)
	} else {
		lastBlock, err := net.Blockchain.GetBlock(net.Blockchain.LastHash)
		if err != nil {
			logrus.Info(err)
		}
		logrus.Info(block.Height)
		valid := block.IsBlockValid(lastBlock)
		logrus.Infof("Block validity: %s", strconv.FormatBool(valid))
		if valid {
			net.Blockchain.AddBlock(block)
		} else {
			CloseDB(net.Blockchain)
			log.Fatalf("We discovered an invalid block of height: %d", block.Height)
		}
	}

	logrus.Printf("Added block %x \n", block.Hash)
	logrus.Info("Block in transit %d", len(blocksInTransit))
	if len(blocksInTransit) > 0 {
		blockHash := blocksInTransit[0]

		net.SendGetData(payload.SendFrom, "block", blockHash)
		blocksInTransit = blocksInTransit[1:]
	} else {
		UTXO := blockchain.UXTOSet{net.Blockchain}
		UTXO.Compute()
	}
}
func (net *Network) SendGetData(peerId string, _type string, id []byte) {
	payload := GobEncode(GetData{net.Host.ID().Pretty(), _type, id})
	request := append(CmdToBytes("getdata"), payload...)
	net.GeneralChannel.Publish("Recieved getdata", request, peerId)
}

func (net *Network) HandleGetData(content *ChannelContent) {
	var buff bytes.Buffer
	var payload GetData

	buff.Write(content.Payload[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)

	if err != nil {
		log.Panic(err)
	}

	if payload.Type == "block" {
		block, err := net.Blockchain.GetBlock([]byte(payload.ID))
		if err != nil {
			return
		}

		net.SendBlock(payload.SendFrom, &block)
	}
}

func (net *Network) SendInv(peerId string, _type string, items [][]byte) {
	inventory := Inv{net.Host.ID().Pretty(), _type, items}
	payload := GobEncode(inventory)
	request := append(CmdToBytes("inv"), payload...)
	net.GeneralChannel.Publish("Recieved inventory", request, peerId)
}

func (net *Network) HandleInv(content *ChannelContent) {
	var buff bytes.Buffer
	var payload Inv

	buff.Write(content.Payload[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)

	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("Recieved inventory with %d %s \n", len(payload.Items), payload.Type)

	if payload.Type == "block" {
		blocksInTransit = payload.Items

		blockHash := payload.Items[0]
		net.SendGetData(payload.SendFrom, "block", blockHash)

		newInTransit := [][]byte{}
		for _, b := range blocksInTransit {
			if bytes.Compare(b, blockHash) != 0 {
				newInTransit = append(newInTransit, b)
			}
		}
		blocksInTransit = newInTransit
	}
}

func (net *Network) SendGetBlocks(peerId string, height int) {
	payload := GobEncode(GetBlocks{net.Host.ID().Pretty(), height})
	request := append(CmdToBytes("getblocks"), payload...)
	net.GeneralChannel.Publish(" Recieved getblocks", request, peerId)
}

func (net *Network) HandleGetBlocks(content *ChannelContent) {
	var buff bytes.Buffer
	var payload GetBlocks

	buff.Write(content.Payload[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)

	if err != nil {
		log.Panic(err)
	}

	blockHashes := net.Blockchain.GetBlockHashes(payload.Height)
	fmt.Println("LENGTH:", len(blockHashes))
	net.SendInv(payload.SendFrom, "block", blockHashes)
}

func (net *Network) SendVersion(peer string) {
	bestHeight := net.Blockchain.GetBestHeight()
	payload := GobEncode(Version{
		version,
		bestHeight,
		net.Host.ID().Pretty(),
	})
	request := append(CmdToBytes("version"), payload...)
	net.GeneralChannel.Publish("Recieved send version", request, peer)
}

func (net *Network) HandleVersion(content *ChannelContent) {
	var buff bytes.Buffer
	var payload Version

	buff.Write(content.Payload[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)

	if err != nil {
		log.Panic(err)
	}

	bestHeight := net.Blockchain.GetBestHeight()
	otherHeight := payload.BestHeight
	fmt.Println(bestHeight, otherHeight)
	if bestHeight < otherHeight {
		net.SendGetBlocks(payload.SendFrom, bestHeight)
	} else if bestHeight > otherHeight {
		net.SendVersion(payload.SendFrom)
	}
}

func StartNode(listenPort, minerAddress string, miner bool) *Network {
	var r io.Reader
	r = rand.Reader
	MinerAddress = minerAddress
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	chain := blockchain.ContinueBlockchain()

	defer chain.Database.Close()
	go CloseDB(chain)

	// Creates a new RSA key pair for this host.
	prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		panic(err)
	}

	transports := libp2p.ChainOptions(
		libp2p.Transport(tcp.NewTCPTransport),
		libp2p.Transport(ws.New),
	)

	muxers := libp2p.ChainOptions(
		libp2p.Muxer("/yamux/1.0.0", yamux.DefaultTransport),
		libp2p.Muxer("/mplex/6.7.0", mplex.DefaultTransport),
	)

	// security := libp2p.Security(secio.ID, secio.New)
	if len(listenPort) == 0 {
		listenPort = "0"
	}

	listenAddrs := libp2p.ListenAddrStrings(
		fmt.Sprintf("/ip4/0.0.0.0/tcp/%s", listenPort),
		fmt.Sprintf("/ip4/0.0.0.0/tcp/%s/ws", listenPort),
	)

	host, err := libp2p.New(
		ctx,
		transports,
		listenAddrs,
		muxers,
		libp2p.Identity(prvKey),
	)
	if err != nil {
		panic(err)
	}
	for _, addr := range host.Addrs() {
		fmt.Println("Listening on", addr)
	}
	log.Info("Host created: ", host.ID())

	// create a new PubSub service using the GossipSub router for general room
	pub, err := pubsub.NewGossipSub(ctx, host)
	if err != nil {
		panic(err)
	}

	GeneralChannel, _ := JoinChannel(ctx, pub, host.ID(), GeneralChannel, true)
	subscribe := false
	if miner {
		subscribe = true
	}
	MiningChannel, _ := JoinChannel(ctx, pub, host.ID(), MiningChannel, subscribe)
	ui := NewCLIUI(GeneralChannel, MiningChannel)

	// setup peer discovery
	err = SetupDiscovery(ctx, host)
	if err != nil {
		panic(err)
	}
	network := &Network{
		Host:           host,
		GeneralChannel: GeneralChannel,
		MiningChannel:  MiningChannel,
		Blockchain:     chain,
	}
	err = RequestBlocks(network)
	if err != nil {
		panic(err)
	}
	if err = ui.Run(network); err != nil {
		printErr("error running text UI: %s", err)
	}

	return network
}

func RequestBlocks(net *Network) error {
	peers := net.GeneralChannel.ListPeers()
	// Send version
	if len(peers) > 0 {
		net.SendVersion(peers[0].Pretty())
	}
	return nil
}
func SetupDiscovery(ctx context.Context, host host.Host) error {

	// Start a DHT, for use in peer discovery. We can't just make a new DHT
	// client because we want each peer to maintain its own local copy of the
	// DHT, so that the bootstrapping node of the DHT can go down without
	// inhibiting future peer discovery.
	kademliaDHT, err := dht.New(ctx, host)
	if err != nil {
		panic(err)
	}

	// Bootstrap the DHT. In the default configuration, this spawns a Background
	// thread that will refresh the peer table every five minutes.
	log.Info("Bootstrapping the DHT")
	if err = kademliaDHT.Bootstrap(ctx); err != nil {
		panic(err)
	}

	// Let's connect to the bootstrap nodes first. They will tell us about the
	// other nodes in the network.

	var wg sync.WaitGroup
	for _, peerAddr := range dht.DefaultBootstrapPeers {
		peerinfo, _ := peer.AddrInfoFromP2pAddr(peerAddr)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := host.Connect(ctx, *peerinfo); err != nil {
				log.Error(err)
			} else {
				log.Info("Connection established with bootstrap node:", *peerinfo)
			}
		}()
	}
	wg.Wait()

	// We use a rendezvous point "meet me here" to announce our location.
	// This is like telling your friends to meet you at the Eiffel Tower.
	log.Info("Announcing ourselves...")
	routingDiscovery := discovery.NewRoutingDiscovery(kademliaDHT)
	discovery.Advertise(ctx, routingDiscovery, "rendezvous")
	log.Info("Successfully announced!")

	// Now, look for others who have announced
	// This is like your friend telling you the location to meet you.
	log.Info("Searching for other peers...")
	peerChan, err := routingDiscovery.FindPeers(ctx, "rendezvous")
	if err != nil {
		panic(err)
	}

	// Finally we open streams to the newly discovered peers.
	for peer := range peerChan {
		if peer.ID == host.ID() {
			continue
		}
		log.Debug("Found peer:", peer)

		log.Debug("Connecting to:", peer)
		err := host.Connect(context.Background(), peer)
		if err != nil {
			log.Warningf("Error connecting to peer %s: %s\n", peer.ID.Pretty(), err)
			continue
		}
		log.Info("Connected to:", peer)
	}

	return nil
}

// printErr is like fmt.Printf, but writes to stderr.
func printErr(m string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, m, args...)
}
