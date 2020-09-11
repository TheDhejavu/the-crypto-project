package blockchain

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"time"
)

type Block struct {
	Timestamp    int64
	Hash         []byte
	PrevHash     []byte
	Transactions []*Transaction
	Nonce        int
	Height       int
	MerkleRoot   []byte
	Difficulty   int
}

// Use Merkle Tree to hash Transactions
func (block *Block) HashTransactions() []byte {
	var txHashes [][]byte

	for _, tx := range block.Transactions {
		txHashes = append(txHashes, tx.Serializer())
	}

	tree := NewMerkleTree(txHashes)
	return tree.RootNode.Data
}

func CreateBlock(txs []*Transaction, prevHash []byte, height int) *Block {
	block := &Block{
		time.Now().Unix(),
		[]byte{},
		prevHash,
		txs,
		0,
		height,
		[]byte{},
		Difficulty,
	}
	pow := NewProof(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce
	//Set MerkleRoot
	block.MerkleRoot = block.HashTransactions()

	return block
}

// Genesis block
func Genesis(MinerTx *Transaction) *Block {
	return CreateBlock([]*Transaction{MinerTx}, []byte{}, 0)
}

// Util function for serializing blockchain data
func (b *Block) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	err := encoder.Encode(b)
	Handle(err)
	return res.Bytes()
}

// Util function for De-serializing blockchain data
func DeSerialize(data []byte) *Block {
	var block Block
	encoder := gob.NewDecoder(bytes.NewReader(data))

	err := encoder.Decode(&block)
	Handle(err)
	return &block
}

// Check if the block is valid by confirming variety of information
// in the block
func IsBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Height+1 != newBlock.Height {
		return false
	}
	res := bytes.Compare(oldBlock.Hash, newBlock.PrevHash)
	if res != 0 {
		return false
	}

	return true
}

func ConstructJSON(buffer *bytes.Buffer, block *Block) {
	buffer.WriteString("{")
	buffer.WriteString(fmt.Sprintf("\"%s\":\"%d\",", "Timestamp", block.Timestamp))
	buffer.WriteString(fmt.Sprintf("\"%s\":\"%x\",", "PrevHash", block.PrevHash))

	buffer.WriteString(fmt.Sprintf("\"%s\":\"%x\",", "Hash", block.Hash))

	buffer.WriteString(fmt.Sprintf("\"%s\":%d,", "Difficulty", block.Difficulty))

	buffer.WriteString(fmt.Sprintf("\"%s\":%d,", "Nonce", block.Nonce))

	buffer.WriteString(fmt.Sprintf("\"%s\":\"%x\"", "MerkleRoot", block.MerkleRoot))
	buffer.WriteString("}")
}

func (bs *Block) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("[")
	ConstructJSON(buffer, bs)
	buffer.WriteString("]")
	return buffer.Bytes(), nil
}
