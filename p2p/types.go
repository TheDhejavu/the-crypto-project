package p2p

import (
	"github.com/libp2p/go-libp2p-core/host"
	blockchain "github.com/workspace/the-crypto-project/core"
)

type Network struct {
	Host             host.Host
	GeneralChannel   *Channel
	MiningChannel    *Channel
	FullNodesChannel *Channel
	Blockchain       *blockchain.Blockchain
	Blocks           chan *blockchain.Block
	Transactions     chan *blockchain.Transaction
	Miner            bool
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
type Tx struct {
	SendFrom    string
	Transaction []byte
}
type Block struct {
	SendFrom string
	Block    []byte
}

type TxFromPool struct {
	SendFrom string
	Count    int
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
