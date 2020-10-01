package memopool

import (
	"encoding/hex"

	blockchain "github.com/workspace/the-crypto-project/core"
)

type MemoPool struct {
	Pending map[string]blockchain.Transaction
	Queued  map[string]blockchain.Transaction
}

func (memo *MemoPool) Move(tnx blockchain.Transaction, to string) {
	if to == "pending" {
		memo.Remove(hex.EncodeToString(tnx.ID), "queued")
		memo.Pending[hex.EncodeToString(tnx.ID)] = tnx
	}

	if to == "queued" {
		memo.Remove(hex.EncodeToString(tnx.ID), "pending")
		memo.Pending[hex.EncodeToString(tnx.ID)] = tnx
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


func (memo *MemoPool) RemoveFromAll(txID string) {
	
	delete(memo.Queued, txID)
	delete(memo.Pending, txID)

}