package memopool

import (
	"encoding/hex"
	"sync"

	blockchain "github.com/workspace/the-crypto-project/core"
)

// Memory pool Data-structure

type MemoPool struct {
	Pending map[string]blockchain.Transaction
	Queued  map[string]blockchain.Transaction
	Wg      sync.WaitGroup
}

func (memo *MemoPool) Move(tnx blockchain.Transaction, to string) {
	if to == "pending" {
		memo.Remove(hex.EncodeToString(tnx.ID), "queued")
		memo.Pending[hex.EncodeToString(tnx.ID)] = tnx
	}

	if to == "queued" {
		memo.Remove(hex.EncodeToString(tnx.ID), "pending")
		memo.Queued[hex.EncodeToString(tnx.ID)] = tnx
	}
}

func (memo *MemoPool) Add(tnx blockchain.Transaction) {
	memo.Pending[hex.EncodeToString(tnx.ID)] = tnx
}

func (memo *MemoPool) Remove(txID string, from string) {
	if from == "queued" {
		delete(memo.Queued, txID)
		return
	}

	if from == "pending" {
		delete(memo.Pending, txID)
		return
	}
}

func (memo *MemoPool) GetPendingTransactions(count int) (txs [][]byte) {
	i := 0
	for _, tx := range memo.Pending {
		txs = append(txs, tx.ID)
		if i == count {
			break
		}
		i++
	}
	return txs
}

func (memo *MemoPool) RemoveFromAll(txID string) {
	delete(memo.Queued, txID)
	delete(memo.Pending, txID)
}

func (memo *MemoPool) ClearAll() {
	memo.Pending = map[string]blockchain.Transaction{}
	memo.Queued = map[string]blockchain.Transaction{}
}
