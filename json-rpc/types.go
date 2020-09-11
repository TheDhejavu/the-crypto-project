package rpc

import (
	"bytes"

	blockchain "github.com/workspace/the-crypto-project/core"
)

type Args struct {
	Address string
}

type SendArgs struct {
	SendFrom string
	SendTo   string
	Amount   float64
	Mine     bool
}

type BlockArgs struct {
	Address string
	Height  int
}

type Blocks []*blockchain.Block

func (bs *Blocks) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("[")
	count := len(*bs)
	for _, block := range *bs {
		blockchain.ConstructJSON(buffer, block)
		count -= 1
		if count != 0 {
			buffer.WriteString(",")
		}
	}
	buffer.WriteString("]")
	return buffer.Bytes(), nil
}
