package utils

import (
	"os"
	"runtime"
	"syscall"

	blockchain "github.com/workspace/the-crypto-project/core"
	"gopkg.in/vrecan/death.v3"
)

func CloseDB(chain *blockchain.Blockchain) {
	d := death.NewDeath(syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	d.WaitForDeathWithFunc(func() {
		defer os.Exit(1)
		defer runtime.Goexit()
		chain.Database.Close()
	})
}
