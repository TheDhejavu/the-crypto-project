package blockchain

import badger "github.com/dgraph-io/badger"

type BlockchainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

func (chain *Blockchain) Iterator() *BlockchainIterator {
	if chain.LastHash == nil {
		return nil
	}
	return &BlockchainIterator{chain.LastHash, chain.Database}
}

func (iter *BlockchainIterator) Next() *Block {
	var block *Block
	var encodedBlock []byte

	//Read
	err := iter.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.CurrentHash)
		Handle(err)
		encodedBlock, err = item.ValueCopy(nil)
		block = DeSerialize(encodedBlock)
		return err
	})
	Handle(err)

	iter.CurrentHash = block.PrevHash
	return block
}
