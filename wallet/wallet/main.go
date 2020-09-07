package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/workspace/the-crypto-project/wallet"
)

const (
	cwd = true
)

func PrintWalletAddress(address string, w wallet.Wallet) {
	var lines []string
	lines = append(lines, fmt.Sprintf("======ADDRESS:======\n %s ", address))
	lines = append(lines, fmt.Sprintf("======PUBLIC KEY:======\n %x", w.PublicKey))
	lines = append(lines, fmt.Sprintf("======PRIVATE KEY:======\n %x", w.PrivateKey.D.Bytes()))
	fmt.Println(strings.Join(lines, "\n"))
}

func main() {
	var cmdGenerate = &cobra.Command{
		Use:   "generate",
		Short: "Generate new wallet and print",
		Run: func(cmd *cobra.Command, args []string) {

			wallets, _ := wallet.InitializeWallets(cwd)
			address := wallets.AddWallet()
			currentPath := true
			wallets.SaveFile(currentPath)
			w := wallets.GetWallet(address)
			PrintWalletAddress(address, w)
		},
	}
	var Address string
	var cmdPrint = &cobra.Command{
		Use:   "print",
		Short: "Print wallet address",
		Run: func(cmd *cobra.Command, args []string) {
			var w wallet.Wallet
			var address string
			wallets, _ := wallet.InitializeWallets(cwd)
			if Address != "" {
				if !wallet.ValidateAddres(Address) {
					log.Panic("Invalid address")
				}
				w = wallets.GetWallet(Address)
				PrintWalletAddress(Address, w)
			} else {
				count := 1
				for address = range wallets.Wallets {
					w = *wallets.Wallets[address]
					fmt.Println("")
					PrintWalletAddress(address, w)
					count++
				}
			}
		},
	}
	cmdPrint.PersistentFlags().StringVar(&Address, "address", "", "Wallet address")

	var rootCmd = &cobra.Command{Use: "wallet"}
	rootCmd.AddCommand(cmdGenerate, cmdPrint)
	rootCmd.Execute()
}
