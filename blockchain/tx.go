package blockchain

import (
	"bytes"
	"encoding/gob"

	"github.com/workspace/go_blockchain/wallet"
)

// Input represents debit
type TxInput struct {
	ID        []byte
	Out       int
	Signature []byte
	PubKey    []byte
}

// output represents credit
type TxOutputs struct {
	Outputs []TxOutput
}

// output represents credit
type TxOutput struct {
	Value      int
	PubKeyHash []byte
}

func NewTXOutput(value int, address string) *TxOutput {
	txo := &TxOutput{value, nil}
	txo.Lock([]byte(address))

	return txo
}

func (in *TxInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := wallet.PublicKeyHash(in.PubKey)

	return bytes.Compare(lockingHash, pubKeyHash) == 0
}

func (out *TxOutput) Lock(address []byte) {
	pubKeyHash := wallet.Base58Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]

	out.PubKeyHash = pubKeyHash
}

func (out *TxOutput) IsLockWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}


func (outputs *TxOutputs) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	err := encoder.Encode(outputs)
	Handle(err)
	return res.Bytes()
}

func DeSerializeOutputs(data []byte) TxOutputs {
	var outputs TxOutputs
	encoder := gob.NewDecoder(bytes.NewReader(data))

	err := encoder.Decode(&outputs)
	Handle(err)
	return outputs
}
