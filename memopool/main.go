package memopool

import (
	blockchain "github.com/workspace/the-crypto-project/core"
)

type MemoPool struct {
	Local  map[string]blockchain.Transaction
	Global map[string]blockchain.Transaction
}

func (memo *MemoPool) Get() bool {
	return false
}
