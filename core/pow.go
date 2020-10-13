package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"math"
	"math/big"

	log "github.com/sirupsen/logrus"
)

const Difficulty = 5

type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

// /https://imil.net/blog/posts/2019/proof-of-work-based-blockchain-explained-with-golang/
// Create a new Proof.
func NewProof(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-Difficulty))

	pow := &ProofOfWork{b, target}
	log.Infof("Target: %x\n", target)

	return pow
}

// Initialize the block data by concatenating
// the transaction hash + prevHash + nonce + POW Difficulty
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

// Execute the Proof Of Work by incrementing the nonce
// util the  hash falls below the the target value base on the Difficulty level
func (pow *ProofOfWork) Run() (int, []byte) {
	var initHash big.Int
	var hash [32]byte

	nonce := 0

	for nonce = 0; nonce < math.MaxInt64; nonce++ {
		info := pow.InitData(nonce)
		hash = sha256.Sum256(info)

		log.Infof("Pow: \r%x", hash)
		initHash.SetBytes(hash[:])

		if initHash.Cmp(pow.Target) == -1 {
			log.Info("Found!")
			break
		}
	}
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
