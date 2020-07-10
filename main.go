package main

import (
	"os"

	"github.com/workspace/go_blockchain/cli"
)

func main() {
	defer os.Exit(0)
	cli := &cli.CommandLine{}
	cli.Run()
}
