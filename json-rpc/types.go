package rpc

import (
	"bytes"
	"fmt"

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
type Blocks []*blockchain.Block

func (bs *Blocks) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("[")
	count := len(*bs)
	for _, block := range *bs {
		buffer.WriteString("{")
		buffer.WriteString(fmt.Sprintf("\"%s\":\"%d\",", "Timestamp", block.Timestamp))
		buffer.WriteString(fmt.Sprintf("\"%s\":\"%x\",", "PrevHash", block.PrevHash))

		buffer.WriteString(fmt.Sprintf("\"%s\":\"%x\",", "Hash", block.Hash))

		buffer.WriteString(fmt.Sprintf("\"%s\":%d,", "Difficulty", block.Difficulty))

		buffer.WriteString(fmt.Sprintf("\"%s\":%d,", "Nonce", block.Nonce))

		buffer.WriteString(fmt.Sprintf("\"%s\":\"%x\"", "MerkleRoot", block.MerkleRoot))
		buffer.WriteString("}")
		count -= 1
		if count != 0 {
			buffer.WriteString(",")
		}
	}
	buffer.WriteString("]")
	return buffer.Bytes(), nil
}
