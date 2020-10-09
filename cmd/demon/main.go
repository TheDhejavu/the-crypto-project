package main

import (
	"fmt"
	"os"
	"time"

	"github.com/mattn/go-colorable"
	log "github.com/sirupsen/logrus"
	"github.com/snowzach/rotatefilehook"
	"github.com/spf13/cobra"
	"github.com/workspace/the-crypto-project/cmd/utils"
	blockchain "github.com/workspace/the-crypto-project/core"
	jsonrpc "github.com/workspace/the-crypto-project/json-rpc"
	"github.com/workspace/the-crypto-project/p2p"
	"github.com/workspace/the-crypto-project/util/env"
)

func init() {
	var logLevel = log.InfoLevel

	rotateFileHook, err := rotatefilehook.NewRotateFileHook(rotatefilehook.RotateFileConfig{
		Filename:   "../../logs/console.log",
		MaxSize:    50, // megabytes
		MaxBackups: 3,
		MaxAge:     28, //days
		Level:      logLevel,
		Formatter: &log.JSONFormatter{
			TimestampFormat: time.RFC822,
		},
	})

	if err != nil {
		log.Fatalf("Failed to initialize file rotate hook: %v", err)
	}

	log.SetLevel(logLevel)
	log.SetOutput(colorable.NewColorableStdout())
	log.SetFormatter(&log.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: time.RFC822,
	})
	log.AddHook(rotateFileHook)
}

func main() {
	defer os.Exit(0)
	var conf = env.New()
	var address string
	var instanceId string

	var rpcPort string
	var rpcAddr string
	var rpc bool

	cli := utils.CommandLine{
		Blockchain: &blockchain.Blockchain{
			Database:   nil,
			InstanceId: instanceId,
		},
		P2p: nil,
	}

	/*
	* INIT COMMAND
	 */
	var initCmd = &cobra.Command{
		Use:   "init",
		Short: "Initialize the blockchain and create the genesis block",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {

			cli := cli.UpdateInstance(instanceId, true)
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
			cli := cli.UpdateInstance(instanceId, true)
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
			cli := cli.UpdateInstance(instanceId, true)
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
			cli := cli.UpdateInstance(instanceId, true)
			cli.PrintBlockchain()
		},
	}

	/*
	* NODE COMMAND
	 */
	var minerAddress string
	var miner bool
	var fullNode bool
	var listenPort string
	var nodeCmd = &cobra.Command{
		Use:   "startnode",
		Short: "start a node",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(miner)
			if miner && len(minerAddress) == 0 {
				log.Fatalln("Miner address is required --address")
			}

			cli := cli.UpdateInstance(instanceId, false)
			cli.StartNode(listenPort, minerAddress, miner, fullNode, func(net *p2p.Network) {
				if rpc {
					cli.P2p = net
					go jsonrpc.StartServer(cli, rpc, rpcPort, rpcAddr)
				}
			})
		},
	}
	nodeCmd.Flags().StringVar(&listenPort, "port", conf.ListenPort, "Node listening port")
	nodeCmd.Flags().StringVar(&minerAddress, "address", conf.MinerAddress, "Set miner address")
	nodeCmd.Flags().BoolVar(&miner, "miner", conf.Miner, "Set as true if you are joining the network as a miner")
	nodeCmd.Flags().BoolVar(&fullNode, "fullnode", conf.FullNode, "Set as true if you are joining the network as a miner")

	/*
	* SEND COMMAND
	 */
	var mine bool
	var sendFrom string
	var sendTo string
	var amount float64

	var sendCmd = &cobra.Command{
		Use:   "send",
		Short: "Send x amount of token to address from local wallet address",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			cli.Send(sendFrom, sendTo, amount, mine)
		},
	}
	sendCmd.Flags().StringVar(&sendFrom, "sendfrom", "", "Sender's wallet address")
	sendCmd.Flags().StringVar(&sendTo, "sendto", "", "Reciever's wallet address")
	sendCmd.Flags().Float64Var(&amount, "amount", float64(0), "Amount of token to send")
	sendCmd.Flags().BoolVar(&mine, "mine", false, "Set if you want your Node to mine the transaction instantly")

	var rootCmd = &cobra.Command{
		Use: "demon",
		Run: func(cmd *cobra.Command, args []string) {
			cli := cli.UpdateInstance(instanceId, false)

			if rpc {
				jsonrpc.StartServer(cli, rpc, rpcPort, rpcAddr)
			}
		},
	}
	rootCmd.PersistentFlags().StringVar(&address, "address", "", "Wallet address")

	/*
	* HTTP FLAGS
	 */
	rootCmd.PersistentFlags().StringVar(&rpcPort, "rpcport", "", " HTTP-RPC server listening port (default: 5000)")
	rootCmd.PersistentFlags().StringVar(&rpcAddr, "rpcaddr", "", "HTTP-RPC server listening interface (default: localhost)")
	rootCmd.PersistentFlags().BoolVar(&rpc, "rpc", false, "Enable the HTTP-RPC server")

	rootCmd.PersistentFlags().StringVar(&instanceId, "instanceid", "", "Blockchain instance")
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
