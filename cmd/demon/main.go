package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/workspace/the-crypto-project/cmd/utils"
	jsonrpc "github.com/workspace/the-crypto-project/json-rpc"
	"github.com/workspace/the-crypto-project/util/env"
)

func main() {
	var conf = env.New()
	defer os.Exit(0)
	cli := utils.CommandLine{}
	var address string
	var ListenPort string
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
			if ListenPort == "" {
				ListenPort = conf.ListenPort
			}
			if minerAddress == "" {
				minerAddress = conf.MinerAddress
			}

			cli.StartNode(ListenPort, address, miner)
		},
	}
	nodeCmd.Flags().StringVar(&ListenPort, "ListenPort", "", "Node ID")
	nodeCmd.Flags().StringVar(&minerAddress, "minerAddress", "", "Set miner address")
	nodeCmd.Flags().BoolVar(&miner, "miner", conf.Miner, "Set as true if you are joining the network as a miner")

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
	sendCmd.Flags().StringVar(&sendFrom, "sendFrom", "", "Sender's wallet address")
	sendCmd.Flags().StringVar(&sendTo, "sendTo", "", "Reciever's wallet address")
	sendCmd.Flags().Float64Var(&amount, "amount", float64(0), "Amount of token to send")
	sendCmd.Flags().BoolVar(&mine, "mine", false, "Set if you want your Node to mine the transaction instantly")

	var rpcPort int
	var rpcAddr string
	var rpc bool
	var rootCmd = &cobra.Command{
		Use: "demon",
		Run: func(cmd *cobra.Command, args []string) {
			if rpc {
				jsonrpc.StartServer(rpc, rpcPort, rpcAddr)
			}
		},
	}
	rootCmd.PersistentFlags().StringVar(&address, "address", "", "Wallet address")

	/*
	* HTTP FLAGS
	 */
	rootCmd.PersistentFlags().IntVar(&rpcPort, "rpcPort", 0, " HTTP-RPC server listening port (default: 1245)")
	rootCmd.PersistentFlags().StringVar(&rpcAddr, "rpcAddr", "", "HTTP-RPC server listening interface (default: localhost)")
	rootCmd.PersistentFlags().BoolVar(&rpc, "rpc", false, "Enable the HTTP-RPC server")
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
