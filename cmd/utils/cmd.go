package utils

import (
	"fmt"
	"log"
	"strconv"

	blockchain "github.com/workspace/the-crypto-project/core"
	"github.com/workspace/the-crypto-project/network"
	"github.com/workspace/the-crypto-project/wallet"
)

type CommandLine struct {
	Blockchain *blockchain.Blockchain
}

func (cli *CommandLine) StartNode(ListenPort, minerAddress string, miner bool) {
	if miner {
		fmt.Printf("Starting Node %s as a MINER\n", ListenPort)
	} else {
		fmt.Printf("Starting Node %s\n", ListenPort)
	}
	if len(minerAddress) > 0 {
		if wallet.ValidateAddres(minerAddress) {
			fmt.Println("Mining is on. Address to receive rewards:", minerAddress)
		} else {
			log.Panic("Wrong Miner Address!")
		}
	}

	network.StartServer(ListenPort, minerAddress)
}

func (cli *CommandLine) Send(from string, to string, amount float64, mineNow bool) string {

	if !wallet.ValidateAddres(from) {
		log.Panic("sendTo address is Invalid ")
	}
	if !wallet.ValidateAddres(to) {
		log.Panic("sendFrom address is Invalid ")
	}

	chain := blockchain.ContinueBlockchain()

	defer chain.Database.Close()
	utxos := blockchain.UXTOSet{chain}
	cwd := false
	wallets, err := wallet.InitializeWallets(cwd)
	if err != nil {
		log.Panic(err)
	}

	wallet := wallets.GetWallet(from)

	tx := blockchain.NewTransaction(&wallet, to, amount, &utxos)

	if mineNow {

		cbTx := blockchain.MinerTx(from, "")
		txs := []*blockchain.Transaction{cbTx, tx}
		block := chain.MineBlock(txs)
		utxos.Update(block)
	} else {
		network.SendTx(network.KnownNodes[0], tx)
		fmt.Println("")
	}
	fmt.Println("Success!")

	return "Successfully sent"
}
func (cli *CommandLine) CreateBlockchain(address string) {
	if !wallet.ValidateAddres(address) {
		log.Panic("Invalid address")
	}

	chain := blockchain.InitBlockchain(address)
	defer chain.Database.Close()

	utxos := blockchain.UXTOSet{chain}
	utxos.Compute()
	fmt.Println("Initialized Blockchain Successfully")
}

func (cli *CommandLine) ComputeUTXOs() {
	chain := blockchain.ContinueBlockchain()
	defer chain.Database.Close()

	utxos := blockchain.UXTOSet{chain}
	utxos.Compute()
	count := utxos.CountTransactions()
	fmt.Printf("Rebuild DONE!!!!, there are %d transactions in the utxos set", count)
}
func (cli *CommandLine) GetBalance(address string) string {
	if !wallet.ValidateAddres(address) {
		log.Panic("Invalid address")
	}
	chain := blockchain.ContinueBlockchain()
	defer chain.Database.Close()
	balance := float64(0)
	publicKeyHash := wallet.Base58Decode([]byte(address))
	publicKeyHash = publicKeyHash[1 : len(publicKeyHash)-4]
	utxos := blockchain.UXTOSet{chain}

	UTXOs := utxos.FindUnSpentTransactions(publicKeyHash)
	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of %s:%f\n", address, balance)

	return fmt.Sprintf("Balance of %s:%f", address, balance)
}

func (cli *CommandLine) CreateWallet() string {
	cwd := false
	wallets, _ := wallet.InitializeWallets(cwd)
	address := wallets.AddWallet()
	wallets.SaveFile(cwd)

	fmt.Println(address)
	return address
}

func (cli *CommandLine) ListAddresses() {
	cwd := false
	wallets, err := wallet.InitializeWallets(cwd)
	if err != nil {
		log.Panic(err)
	}
	addresses := wallets.GetAllAddress()

	for _, address := range addresses {
		fmt.Println(address)
	}
}
func (cli *CommandLine) PrintBlockchain() {
	chain := blockchain.ContinueBlockchain()

	defer chain.Database.Close()
	iter := chain.Iterator()

	for {
		block := iter.Next()
		fmt.Printf("PrevHash: %x\n", block.PrevHash)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Printf("Height: %d\n", block.Height)
		pow := blockchain.NewProof(block)
		validate := pow.Validate()
		fmt.Printf("Valid: %s\n", strconv.FormatBool(validate))
		for _, tx := range block.Transactions {
			fmt.Println(tx)
		}
		fmt.Println()

		if len(block.PrevHash) == 0 {
			break
		}
	}
}

func (cli *CommandLine) GetBlockchainData() []*blockchain.Block {
	var blocks []*blockchain.Block
	chain := blockchain.ContinueBlockchain()

	defer chain.Database.Close()
	iter := chain.Iterator()

	for {
		block := iter.Next()
		blocks = append(blocks, block)

		if len(block.PrevHash) == 0 {
			break
		}
	}

	return blocks
}
