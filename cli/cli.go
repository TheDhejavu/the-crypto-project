package cli

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"

	"github.com/workspace/go_blockchain/blockchain"
	"github.com/workspace/go_blockchain/network"
	"github.com/workspace/go_blockchain/wallet"
)

type CommandLine struct {
	Blockchain *blockchain.Blockchain
}

func (cli *CommandLine) PrintUsage() {
	fmt.Println("USAGE:")
	fmt.Println("getbalance -address ADDRESS - get the address balance")
	fmt.Println("createblockchain -address ADDRESS creates the genesis block")
	fmt.Println("printchain - Prints the blocks in the chain")
	fmt.Println("send -from FROM -to To -amount AMOUNT -mine - Send amount to address from your address")
	fmt.Println("createwallet - Creates a new wallet")
	fmt.Println("listaddresses - list out all available adresses")
	fmt.Println("reindexutxo - rebuild Unspent transaction outputs")
	fmt.Println("startnode -miner ADDRESS - start a node with ID specified i NODE_ID env.")
}

func (cli *CommandLine) StartNode(nodeId, minerAddress string) {
	fmt.Printf("Starting Node %s\n", nodeId)
	if len(minerAddress) > 0 {
		if wallet.ValidateAddres(minerAddress) {
			fmt.Println("Mining is on. Address to receive rewards:", minerAddress)
		} else {
			log.Panic("Wrong Miner Address!")
		}
	}

	network.StartServer(nodeId, minerAddress)
}

func (cli *CommandLine) ValidateArgs() {
	if len(os.Args) < 2 {
		cli.PrintUsage()
		// allows us to close the program database
		// gracefully without corrupting our data
		runtime.Goexit()
	}
}

func (cli *CommandLine) Send(from string, to string, amount int, nodeId string, mineNow bool) {
	if !wallet.ValidateAddres(from) {
		log.Panic("Invalid to address")
	}
	if !wallet.ValidateAddres(to) {
		log.Panic("Invalid from address")
	}

	chain := blockchain.ContinueBlockchain(nodeId)
	defer chain.Database.Close()
	utxos := blockchain.UXTOSet{chain}

	wallets, err := wallet.InitializeWallets()
	if err != nil {
		log.Panic(err)     
	}

	wallet := wallets.GetWallet(from)

	tx := blockchain.NewTransaction(&wallet, to, amount, &utxos)

	if mineNow {

		cbTx := blockchain.CoinBaseTx(from, "")
		txs := []*blockchain.Transaction{cbTx, tx}
		block := chain.MineBlock(txs)
		utxos.Update(block)
	} else {
		network.SendTx(network.KnownNodes[0], tx)
		fmt.Println("")
	}
	fmt.Println("Success!")
}
func (cli *CommandLine) CreateBlockchain(address, nodeId string) {
	if !wallet.ValidateAddres(address) {
		log.Panic("Invalid to address")
	}
	chain := blockchain.InitBlockchain(address, nodeId)
	defer chain.Database.Close()

	utxos := blockchain.UXTOSet{chain}
	utxos.ReIndex()
	fmt.Println("INITIALIZED BLOCKCHAIN")
}

func (cli *CommandLine) ReIndexUTXOs(nodeId string) {
	chain := blockchain.ContinueBlockchain(nodeId)
	defer chain.Database.Close()

	utxos := blockchain.UXTOSet{chain}
	utxos.ReIndex()
	count := utxos.CountTransactions()
	fmt.Printf("Rebuild DONE!!!!, there are %d transactions in the utxos set", count)
}
func (cli *CommandLine) GetBalance(address, nodeId string) {
	if !wallet.ValidateAddres(address) {
		log.Panic("Invalid address")
	}
	chain := blockchain.ContinueBlockchain(nodeId)
	defer chain.Database.Close()
	balance := 0
	publicKeyHash := wallet.Base58Decode([]byte(address))
	publicKeyHash = publicKeyHash[1 : len(publicKeyHash)-4]
	utxos := blockchain.UXTOSet{chain}

	UTXOs := utxos.FindUnSpentTransactions(publicKeyHash)
	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of %s:%d\n", address, balance)
}

func (cli *CommandLine) Run() {
	cli.ValidateArgs()

	nodeId := os.Getenv("NODE_ID")

	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	listAddressesCmd := flag.NewFlagSet("listaddresses", flag.ExitOnError)
	reindexutxoCmd := flag.NewFlagSet("reindexutxo", flag.ExitOnError)
	startnodeCmd := flag.NewFlagSet("startnode", flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("address", "", "The address to check balance from")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "")
	sendFrom := sendCmd.String("from", "", "Source waller address")
	sendTo := sendCmd.String("to", "", "Destination of the coin address")
	sendAmount := sendCmd.Int("amount", 0, "amount to send to destination address")
	sendMine := sendCmd.Bool("mine", false, "Mine immediately on the same node")
	startNode := startnodeCmd.String("miner", "", "Enable mining mode and send reward to miners")

	switch os.Args[1] {
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		blockchain.Handle(err)
	case "startnode":
		err := startnodeCmd.Parse(os.Args[2:])
		blockchain.Handle(err)
	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		blockchain.Handle(err)
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		blockchain.Handle(err)
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		blockchain.Handle(err)
	case "reindexutxo":
		err := reindexutxoCmd.Parse(os.Args[2:])
		blockchain.Handle(err)
	case "createwallet":
		err := createWalletCmd.Parse(os.Args[2:])
		blockchain.Handle(err)
	case "listaddresses":
		err := listAddressesCmd.Parse(os.Args[2:])
		blockchain.Handle(err)
	default:
		cli.PrintUsage()
		runtime.Goexit()
	}

	if getBalanceCmd.Parsed() {

		if *getBalanceAddress == "" {
			cli.PrintUsage()
			runtime.Goexit()
		}
		cli.GetBalance(*getBalanceAddress, nodeId)
	}

	if createBlockchainCmd.Parsed() {

		if *createBlockchainAddress == "" {
			cli.PrintUsage()
			runtime.Goexit()
		}
		cli.CreateBlockchain(*createBlockchainAddress, nodeId)
	}

	if sendCmd.Parsed() {

		if *sendFrom == "" || *sendTo == "" {
			cli.PrintUsage()
			runtime.Goexit()
		}
		cli.Send(*sendFrom, *sendTo, *sendAmount, nodeId, *sendMine)
	}

	if createWalletCmd.Parsed() {
		cli.CreateWallet()
	}

	if listAddressesCmd.Parsed() {
		cli.ListAddresses()
	}

	if printChainCmd.Parsed() {
		cli.PrintBlockchain(nodeId)
	}

	if reindexutxoCmd.Parsed() {
		cli.ReIndexUTXOs(nodeId)
	}

	if startnodeCmd.Parsed() {
		nodeId = os.Getenv("NODE_ID")
		cli.StartNode(nodeId, *startNode)
	}

}

func (cli *CommandLine) CreateWallet() {
	wallets, _ := wallet.InitializeWallets()
	address := wallets.AddWallet()
	wallets.SaveFile()

	fmt.Println(address)
}

func (cli *CommandLine) ListAddresses() {
	wallets, _ := wallet.InitializeWallets()
	addresses := wallets.GetAllAddress()

	for _, address := range addresses {
		fmt.Println(address)
	}
}
func (cli *CommandLine) PrintBlockchain(nodeId string) {
	chain := blockchain.ContinueBlockchain(nodeId)

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
