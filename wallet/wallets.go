package wallet

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

const walletsPath = "./tmp/wallets_%s.data"

type Wallets struct {
	Wallets map[string]*Wallet
}

func InitializeWallets(nodeId string) (*Wallets, error) {
	wallets := Wallets{map[string]*Wallet{}}
	err := wallets.LoadFile(nodeId)

	return &wallets, err

}
func (ws *Wallets) GetWallet(address string) Wallet {
	return *ws.Wallets[address]
}
func (ws *Wallets) AddWallet() string {
	wallet := MakeWallet()
	address := fmt.Sprintf("%s", wallet.Address())

	ws.Wallets[address] = wallet

	return address
}
func (ws *Wallets) GetAllAddress() []string {
	var addresses []string
	for address := range ws.Wallets {
		addresses = append(addresses, address)
	}
	return addresses
}
func (ws *Wallets) LoadFile(nodeId string) error {
	walletsFile := fmt.Sprintf(walletsPath, nodeId)
	
	if _, err := os.Stat(walletsFile); os.IsNotExist(err) {
		return err
	}
	var wallets Wallets
	fileContent, err := ioutil.ReadFile(walletsFile)
	if err != nil {
		return err
	}

	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallets)
	if err != nil {
		return err
	}

	ws.Wallets = wallets.Wallets

	return nil
}
func (ws *Wallets) SaveFile(nodeId string) {
	var content bytes.Buffer

	gob.Register(elliptic.P256())

	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(ws)
	if err != nil {
		log.Panic(err)
	}
	walletsFile := fmt.Sprintf(walletsPath, nodeId)
	err = ioutil.WriteFile(walletsFile, content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}
}
