package main

import (
	"fmt"

	"github.com/workspace/go_blockchain/wallet"
)

func main() {
	wallets, _ := wallet.InitializeWallets()
	address := wallets.AddWallet()
	w := wallets.GetWallet(address)
	fmt.Println("ADDRESS:", address)
	fmt.Println("PUBLIC KEY:", w.PublicKey)
	fmt.Println("PRIVATE KEY:", w.PrivateKey)
}
