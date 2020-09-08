package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/workspace/the-crypto-project/blockchain"
	"github.com/workspace/the-crypto-project/network"
	"github.com/workspace/the-crypto-project/util/env"
	"github.com/workspace/the-crypto-project/wallet"
)

type CommandLine struct {
	Blockchain *blockchain.Blockchain
}

func (cli *CommandLine) StartNode(nodeId, minerAddress string, miner bool) {
	if miner {
		fmt.Printf("Starting Node %s as a MINER\n", nodeId)
	} else {
		fmt.Printf("Starting Node %s\n", nodeId)
	}
	if len(minerAddress) > 0 {
		if wallet.ValidateAddres(minerAddress) {
			fmt.Println("Mining is on. Address to receive rewards:", minerAddress)
		} else {
			log.Panic("Wrong Miner Address!")
		}
	}

	network.StartServer(nodeId, minerAddress)
}

func (cli *CommandLine) Send(from string, to string, amount float64, mineNow bool) {

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
}
func (cli *CommandLine) CreateBlockchain(address string) {
	if !wallet.ValidateAddres(address) {
		log.Panic("Invalid address")
	}

	chain := blockchain.InitBlockchain(address)
	defer chain.Database.Close()

	utxos := blockchain.UXTOSet{chain}
	utxos.Compute()
	fmt.Println("INITIALIZED BLOCKCHAIN")
}

func (cli *CommandLine) ComputeUTXOs() {
	chain := blockchain.ContinueBlockchain()
	defer chain.Database.Close()

	utxos := blockchain.UXTOSet{chain}
	utxos.Compute()
	count := utxos.CountTransactions()
	fmt.Printf("Rebuild DONE!!!!, there are %d transactions in the utxos set", count)
}
func (cli *CommandLine) GetBalance(address string) {
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
}

func (cli *CommandLine) CreateWallet() {
	cwd := false
	wallets, _ := wallet.InitializeWallets(cwd)
	address := wallets.AddWallet()
	wallets.SaveFile(cwd)

	fmt.Println(address)
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

func main() {
	var conf = env.New()
	defer os.Exit(0)
	cli := &CommandLine{}
	var address string
	var nodeID string
	/*
	* INIT COMMAND
	 */
	var initCmd = &cobra.Command{
		Use:   "init",
		Short: "Initialize the blockchain and create the genesis block",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			cli.CreateBlockchain(address)
		},
	}

	/*
	* WALLET COMMAND
	 */
	var walletCmd = &cobra.Command{
		Use:   "wallet",
		Short: "Manage wallets",
	}
	var newWalletCmd = &cobra.Command{
		Use:   "new",
		Short: "Create New Wallet",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			cli.CreateWallet()
		},
	}

	var listWalletAddressCmd = &cobra.Command{
		Use:   "listaddress",
		Short: "List out all available adresses",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			cli.ListAddresses()
		},
	}
	var walletBalanceCmd = &cobra.Command{
		Use:   "balance",
		Short: "Get the address balance",
		Run: func(cmd *cobra.Command, args []string) {
			cli.GetBalance(address)
		},
	}
	walletCmd.AddCommand(newWalletCmd, listWalletAddressCmd, walletBalanceCmd)

	/*
	* UTXOS COMMAND
	 */
	var computeutxosCmd = &cobra.Command{
		Use:   "computeutxos",
		Short: "Re-build and Compute Unspent transaction outputs",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			cli.ComputeUTXOs()
		},
	}
	/*
	* PRINT COMMAND
	 */
	var printCmd = &cobra.Command{
		Use:   "print",
		Short: "Print the blocks in the blockchain",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			cli.PrintBlockchain()
		},
	}

	/*
	* NODE COMMAND
	 */
	var minerAddress string
	var miner bool
	var nodeCmd = &cobra.Command{
		Use:   "startnode",
		Short: "start a node",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			if nodeID == "" {
				nodeID = conf.NodeId
			}
			if minerAddress == "" {
				minerAddress = conf.MinerAddress
			}

			cli.StartNode(nodeID, address, miner)
		},
	}
	nodeCmd.Flags().StringVar(&nodeID, "nodeID", "", "Node ID")
	nodeCmd.Flags().StringVar(&minerAddress, "minerAddress", "", "Set miner address")
	nodeCmd.Flags().BoolVar(&miner, "miner", conf.Miner, "Set as true if you are joining the network as a miner")

	/*
	* SEND COMMAND
	 */
	var mine bool
	var sendFrom string
	var sendTo string
	var sendAmount float64

	var sendCmd = &cobra.Command{
		Use:   "send",
		Short: "Send x amount of token to address from local wallet address",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			cli.Send(sendFrom, sendTo, sendAmount, mine)
		},
	}
	sendCmd.Flags().StringVar(&sendFrom, "sendFrom", "", "Sender's wallet address")
	sendCmd.Flags().StringVar(&sendTo, "sendTo", "", "Reciever's wallet address")
	sendCmd.Flags().Float64Var(&sendAmount, "sendAmount", float64(0), "Amount of token to send")
	sendCmd.Flags().BoolVar(&mine, "mine", false, "Set if you want your Node to mine the transaction instantly")

	var rootCmd = &cobra.Command{Use: "chain"}
	rootCmd.PersistentFlags().StringVar(&address, "address", "", "Wallet address")

	rootCmd.AddCommand(
		initCmd,
		walletCmd,
		computeutxosCmd,
		sendCmd,
		printCmd,
		nodeCmd,
	)
	rootCmd.Execute()
}
