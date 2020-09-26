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
	"sync"
	"syscall"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
	discovery "github.com/libp2p/go-libp2p-discovery"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	mplex "github.com/libp2p/go-libp2p-mplex"
	secio "github.com/libp2p/go-libp2p-secio"
	yamux "github.com/libp2p/go-libp2p-yamux"
	tcp "github.com/libp2p/go-tcp-transport"
	ws "github.com/libp2p/go-ws-transport"
	maddr "github.com/multiformats/go-multiaddr"
	log "github.com/sirupsen/logrus"
	"gopkg.in/vrecan/death.v3"

	"github.com/libp2p/go-libp2p-core/host"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	blockchain "github.com/workspace/the-crypto-project/core"
)

type addrList []maddr.Multiaddr

type Version struct {
	Version    int
	BestHeight int
	NodeID     peer.ID
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
	memoryPool      = make(map[string]blockchain.Transaction)
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

func SendVersion(peer peer.ID, channel *Channel, host host.Host, chain *blockchain.Blockchain) {
	bestHeight := chain.GetBestHeight()
	payload := GobEncode(Version{
		version,
		bestHeight,
		host.ID(),
	})
	request := append(CmdToBytes("version"), payload...)
	channel.Publish("hello world", request, ShortID(peer))
}

func HandleVersion(content *ChannelContent, chain *blockchain.Blockchain) {
	var buff bytes.Buffer
	var payload Version

	buff.Write(content.Payload[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)

	if err != nil {
		log.Panic(err)
	}

	bestHeight := chain.GetBestHeight()
	otherHeight := payload.BestHeight
	fmt.Println(bestHeight, otherHeight)
	if bestHeight < otherHeight {
		// SendGetBlocks(payload.AddrFrom)
	} else if bestHeight > otherHeight {
		// SendVersion(payload.AddrFrom, chain)
	}
}

func CloseDB(chain *blockchain.Blockchain) {
	d := death.NewDeath(syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	d.WaitForDeathWithFunc(func() {
		defer os.Exit(1)
		defer runtime.Goexit()
		chain.Database.Close()
	})
}

func StartNode(listenPort, minerAddress string, miner bool) {
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

	security := libp2p.Security(secio.ID, secio.New)
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
		security,
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
	err = RequestBlocks(GeneralChannel, host, chain)
	if err != nil {
		panic(err)
	}
	if err = ui.Run(chain); err != nil {
		printErr("error running text UI: %s", err)
	}
}

func RequestBlocks(generalChannel *Channel, host host.Host, chain *blockchain.Blockchain) error {
	peers := generalChannel.ListPeers()
	// Send version
	if len(peers) > 0 {
		SendVersion(peers[0], generalChannel, host, chain)
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
