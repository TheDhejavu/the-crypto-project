package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/workspace/the-crypto-project/wallet"
)

type Transaction struct {
	ID      []byte
	Inputs  []TxInput
	Outputs []TxOutput
}

func (tx *Transaction) Serializer() []byte {
	var encoded bytes.Buffer

	encode := gob.NewEncoder(&encoded)
	err := encode.Encode(tx)
	Handle(err)

	return encoded.Bytes()
}

func (tx *Transaction) Hash() []byte {
	var hash [32]byte

	txCopy := *tx
	txCopy.ID = []byte{}

	hash = sha256.Sum256(txCopy.Serializer())
	return hash[:]
}

func (tx *Transaction) IsMinerTx() bool {
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].ID) == 0 && tx.Inputs[0].Out == -1
}

func (tx *Transaction) Sign(privKey ecdsa.PrivateKey, prevTXs map[string]Transaction) {

	if tx.IsMinerTx() {
		return
	}

	for _, in := range tx.Inputs {
		if prevTXs[hex.EncodeToString(in.ID)].ID == nil {
			log.Fatal("ERROR: Previous Transaction is not correct")
		}
	}

	txCopy := tx.TrimmedCopy()

	for inId, in := range txCopy.Inputs {
		prevTX := prevTXs[hex.EncodeToString(in.ID)]
		txCopy.Inputs[inId].Signature = nil
		//look for the transaction output that produced this input, then sign it with
		// the rest of the data
		txCopy.Inputs[inId].PubKey = prevTX.Outputs[in.Out].PubKeyHash
		dataToSign := fmt.Sprintf("%x\n", txCopy)

		r, s, err := ecdsa.Sign(rand.Reader, &privKey, []byte(dataToSign))
		Handle(err)
		signature := append(r.Bytes(), s.Bytes()...)

		tx.Inputs[inId].Signature = signature
		txCopy.Inputs[inId].PubKey = nil
	}
}

func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	for _, in := range tx.Inputs {
		inputs = append(inputs, TxInput{in.ID, in.Out, nil, nil})
	}

	for _, out := range tx.Outputs {
		outputs = append(outputs, TxOutput{out.Value, out.PubKeyHash})
	}

	txCopy := Transaction{tx.ID, inputs, outputs}

	return txCopy
}

// Create new Transaction
func NewTransaction(w *wallet.Wallet, to string, amount float64, utxo *UXTOSet) (*Transaction, error) {
	var inputs []TxInput
	var outputs []TxOutput

	publicKeyHash := wallet.PublicKeyHash(w.PublicKey)

	acc, validoutputs := utxo.FindSpendableOutputs(publicKeyHash, amount)
	if acc < amount {
		err := errors.New("You dont have Enough Amount...")
		return nil, err
	}

	from := fmt.Sprintf("%s", w.Address())

	for txId, outs := range validoutputs {
		txID, err := hex.DecodeString(txId)

		Handle(err)
		for _, out := range outs {
			input := TxInput{txID, out, nil, w.PublicKey}
			inputs = append(inputs, input)
		}
	}

	outputs = append(outputs, *NewTXOutput(amount, to))
	if acc > amount {
		outputs = append(outputs, *NewTXOutput(acc-amount, from))
	}

	tx := Transaction{nil, inputs, outputs}
	tx.ID = tx.Hash()

	// Sign the new transaction with wallet Private Key
	utxo.Blockchain.SignTransaction(w.PrivateKey, &tx)

	return &tx, nil
}

func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool {
	if tx.IsMinerTx() {
		return true
	}

	for _, in := range tx.Inputs {
		if prevTXs[hex.EncodeToString(in.ID)].ID == nil {
			log.Fatal("ERROR: Previous Transaction is not valid")
		}
	}

	txCopy := tx.TrimmedCopy()
	curve := elliptic.P256()
	for inId, in := range tx.Inputs {
		prevTX := prevTXs[hex.EncodeToString(in.ID)]
		txCopy.Inputs[inId].Signature = nil
		txCopy.Inputs[inId].PubKey = prevTX.Outputs[in.Out].PubKeyHash

		r := big.Int{}
		s := big.Int{}
		sigLen := len(in.Signature)
		r.SetBytes(in.Signature[:(sigLen / 2)])
		s.SetBytes(in.Signature[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(in.PubKey)
		x.SetBytes(in.PubKey[:(keyLen / 2)])
		y.SetBytes(in.PubKey[(keyLen / 2):])

		dataToVerify := fmt.Sprintf("%x\n", txCopy)

		rawPubKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}
		if ecdsa.Verify(&rawPubKey, []byte(dataToVerify), &r, &s) == false {
			return false
		}
		txCopy.Inputs[inId].PubKey = nil
	}

	return true
}

// Helper function for displaying transaction data in the console
func (tx *Transaction) String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("---Transaction: %x", tx.ID))

	for i, input := range tx.Inputs {
		lines = append(lines, fmt.Sprintf("	Input (%d):", i))
		lines = append(lines, fmt.Sprintf(" 	 	TXID: %x", input.ID))
		lines = append(lines, fmt.Sprintf("		Out: %d", input.Out))
		lines = append(lines, fmt.Sprintf(" 	 	Signature: %x", input.Signature))
		lines = append(lines, fmt.Sprintf("		PubKey: %x", input.PubKey))
	}

	for i, out := range tx.Outputs {
		lines = append(lines, fmt.Sprintf("	Output (%d):", i))
		lines = append(lines, fmt.Sprintf(" 	 	Value: %f", out.Value))
		lines = append(lines, fmt.Sprintf("		PubkeyHash: %x", out.PubKeyHash))
	}

	return strings.Join(lines, "\n")
}

// Miner Transaction with Input && Output credited with 20.000 token for the workdone
// No Signature is required for the miner transaction Input
func MinerTx(to, data string) *Transaction {
	if data == "" {
		randData := make([]byte, 24)
		_, err := rand.Read(randData)
		Handle(err)
		data = fmt.Sprintf("%x", randData)
	}

	txIn := TxInput{[]byte{}, -1, nil, []byte(data)}
	txOut := NewTXOutput(20.000, to)

	tx := Transaction{nil, []TxInput{txIn}, []TxOutput{*txOut}}

	tx.ID = tx.Hash()

	return &tx
}
