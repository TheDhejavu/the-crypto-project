package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
)

const Difficulty = 5

type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

func NewProof(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-Difficulty))

	pow := &ProofOfWork{b, target}

	return pow
}

func (pow *ProofOfWork) InitData(nonce int) []byte {
	info := bytes.Join(
		[][]byte{
			pow.Block.HashTransactions(),
			pow.Block.PrevHash,
			ToByte(int64(nonce)),
			ToByte(int64(Difficulty)),
		}, []byte{})

	return info
}

func (pow *ProofOfWork) Run() (int, []byte) {
	var initHash big.Int
	var hash [32]byte

	nonce := 0

	for nonce < math.MaxInt64 {
		info := pow.InitData(nonce)
		hash = sha256.Sum256(info)

		fmt.Printf("\r%x", hash)
		initHash.SetBytes(hash[:])
		
		if initHash.Cmp(pow.Target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Println()
	return nonce, hash[:]
}

func (pow *ProofOfWork) Validate() bool {
	var initHash big.Int
	var hash [32]byte

	info := pow.InitData(pow.Block.Nonce)
	hash = sha256.Sum256(info)

	initHash.SetBytes(hash[:])

	return initHash.Cmp(pow.Target) == -1
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func ToByte(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	check(err)

	return buff.Bytes()
}
