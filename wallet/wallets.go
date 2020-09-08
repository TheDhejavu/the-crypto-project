package wallet

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
)

var (
	_, b, _, _ = runtime.Caller(0)

	// Root folder of this project
	Root            = filepath.Join(filepath.Dir(b), "../")
	walletsPath     = path.Join(Root, "/tmp/")
	walletsFilename = "wallets.data"
)

type Wallets struct {
	Wallets map[string]*Wallet
}

func InitializeWallets(cwd bool) (*Wallets, error) {
	wallets := Wallets{map[string]*Wallet{}}
	err := wallets.LoadFile(cwd)

	return &wallets, err
}
func (ws *Wallets) GetWallet(address string) Wallet {
	var wallet *Wallet
	var ok bool
	w := *ws
	if wallet, ok = w.Wallets[address]; !ok {
		log.Fatalf("Address does not exist")
	}
	return *wallet
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
func (ws *Wallets) LoadFile(cwd bool) error {
	walletsFile := path.Join(walletsPath, walletsFilename)

	if cwd {
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		walletsFile = path.Join(dir, walletsFilename)
	}

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
func (ws *Wallets) SaveFile(cwd bool) {
	walletsFile := path.Join(walletsPath, walletsFilename)

	if cwd {
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		walletsFile = path.Join(dir, walletsFilename)
	}
	var content bytes.Buffer

	gob.Register(elliptic.P256())

	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(ws)
	if err != nil {
		log.Panic(err)
	}

	err = ioutil.WriteFile(walletsFile, content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}
}
