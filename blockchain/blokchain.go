package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	badger "github.com/dgraph-io/badger"
)

// Blockchain struct such that lastHash represents the lastblock hash
// on the ledger
type Blockchain struct {
	LastHash []byte
	Database *badger.DB
}

const (
	dbPath      = "./tmp/blocks"
	dbFile      = "./tmp/blocks/MANIFEST"
	genesisData = "This is the genesis data"
)

// Check if Blockchain Database already exist
func DBExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func ContinueBlockchain() *Blockchain {
	path := dbPath

	if DBExists(path) == false {
		fmt.Println("No Existing Blockchian DB found, create one!")
		runtime.Goexit()
	}
	var lastHash []byte

	opts := badger.DefaultOptions(path)
	db, err := badger.Open(opts)

	Handle(err)

	//Read-Write Operations
	err = db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))

		Handle(err)
		lastHash, err = item.ValueCopy(nil)

		return err
	})

	return &Blockchain{lastHash, db}
}

// Initialize the blockchain by creating the blockchain database
// with a genesis block with an address
func InitBlockchain(address string) *Blockchain {
	var lastHash []byte
	path := dbPath

	if DBExists(path) {
		fmt.Println("Blockchain already exist")
		runtime.Goexit()
	}
	// Open the Badger database located in the /tmp/blocks directory.
	// It will be created if it doesn't exist.
	opts := badger.DefaultOptions(path)
	opts.ValueDir = path
	db, err := OpenDB(path, opts)
	Handle(err)

	//Read-Write Operations
	err = db.Update(func(txn *badger.Txn) error {
		cbtx := MinerTx(address, genesisData)
		fmt.Println("No existing blockchain found")
		genesis := Genesis(cbtx)
		err = txn.Set(genesis.Hash, genesis.Serialize())
		Handle(err)
		err = txn.Set([]byte("lh"), genesis.Hash)
		lastHash = genesis.Hash

		return err
	})
	Handle(err)

	return &Blockchain{lastHash, db}
}

// Add a block to the blockchain
//https://github.com/dgraph-io/badger#read-write-transactions
func (chain *Blockchain) AddBlock(block *Block) *Block {

	//Read-Write Operations
	err := chain.Database.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get(block.Hash); err == nil {
			return nil
		}

		blockData := block.Serialize()
		err := txn.Set(block.Hash, blockData)
		Handle(err)

		// get the last block
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		lastHash, _ := item.ValueCopy(nil)
		item, err = txn.Get(lastHash)
		Handle(err)
		lastBlockData, _ := item.ValueCopy(nil)
		lastBlock := DeSerialize(lastBlockData)

		//check if the current block height is
		// greater than the lastBlock Height
		if block.Height > lastBlock.Height {
			err := txn.Set([]byte("lh"), block.Hash)
			Handle(err)
			chain.LastHash = block.Hash
		}

		return nil
	})

	Handle(err)

	return block
}

// Get Block from the blockchain
func (chain *Blockchain) GetBlock(blockHash []byte) (Block, error) {
	var block Block
	//Read Operations
	err := chain.Database.View(func(txn *badger.Txn) error {
		if item, err := txn.Get(blockHash); err != nil {
			return errors.New("Block does not exist")
		} else {
			blockData, _ := item.ValueCopy(nil)
			// Deserialize the block data
			block = *DeSerialize(blockData)
		}
		return nil
	})

	if err != nil {
		return block, err
	}
	return block, nil
}

//Aggregate and get all block hashes in the blockchain
func (chain *Blockchain) GetBlockHashes() [][]byte {
	var blocks [][]byte

	iter := chain.Iterator()
	for {
		block := iter.Next()
		prevHash := block.PrevHash

		blocks = append(blocks, block.Hash)

		if prevHash == nil {
			break
		}
	}

	return blocks

}

// Get Best height basically gets the height(Index) of the lastBlock
func (chain *Blockchain) GetBestHeight() int {
	var lastBlock Block

	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		lastHash, _ := item.ValueCopy(nil)

		item, err = txn.Get(lastHash)
		Handle(err)
		lastBlockData, _ := item.ValueCopy(nil)
		lastBlock = *DeSerialize(lastBlockData)

		return nil
	})

	Handle(err)

	return lastBlock.Height

}

//Mine Block Creates a new block and adds it to the blockchain
func (chain *Blockchain) MineBlock(transactions []*Transaction) *Block {
	var lastHash []byte
	var lastHeight int

	for _, tx := range transactions {
		if chain.VerifyTransaction(tx) != true {
			log.Panic("Invalid Transaction")
		}
	}
	lastHash = chain.LastHash
	//Populate lastHeight
	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		lastHash, err = item.ValueCopy(nil)

		item, err = txn.Get(lastHash)
		Handle(err)
		lastBlockData, _ := item.ValueCopy(nil)

		lastBlock := DeSerialize(lastBlockData)

		lastHeight = lastBlock.Height
		return err
	})

	Handle(err)

	block := CreateBlock(transactions, lastHash, lastHeight+1)
	// Read-write
	err = chain.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(block.Hash, block.Serialize())
		Handle(err)
		err = txn.Set([]byte("lh"), block.Hash)

		chain.LastHash = lastHash

		return err
	})

	Handle(err)
	return block
}

func DeserializeTransaction(data []byte) Transaction {
	var transaction Transaction

	dec := gob.NewDecoder(bytes.NewReader(data))
	err := dec.Decode(&transaction)
	Handle(err)
	return transaction
}

// Aggregate all Unspent Transaction output from the blockchain
func (chain *Blockchain) FindUTXO() map[string]TxOutputs {
	UTXOs := make(map[string]TxOutputs)
	spentTXOs := make(map[string][]int)

	iter := chain.Iterator()

	for {
		block := iter.Next()

		for _, tx := range block.Transactions {
			//Convert transaction ID to string
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Outputs {
				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}
				//Add to UTXO 
				outs := UTXOs[txID]
				outs.Outputs = append(outs.Outputs, out)
				UTXOs[txID] = outs
			}
			if !tx.IsMinerTx() {
				//Keep Track of Spent Transaction Outputs 
				for _, in := range tx.Inputs {
					inTxID := hex.EncodeToString(in.ID)
					spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Out)
				}
			}
		}
		if len(block.PrevHash) == 0 {
			break
		}
	}
	return UTXOs
}

func (chain *Blockchain) FindTransaction(ID []byte) (Transaction, error) {
	iter := chain.Iterator()

	for {
		block := iter.Next()

		for _, tx := range block.Transactions {
			if bytes.Compare(tx.ID, ID) == 0 {
				return *tx, nil
			}
		}
		if len(block.PrevHash) == 0 {
			break
		}
	}

	return Transaction{}, errors.New("No transaction with id")
}

func (chain *Blockchain) SignTransaction(privKey ecdsa.PrivateKey, tx *Transaction) {
	prevTXs := make(map[string]Transaction)

	for _, in := range tx.Inputs {
		prevTX, err := chain.FindTransaction(in.ID)
		Handle(err)
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}

	tx.Sign(privKey, prevTXs)
}

func (bc *Blockchain) VerifyTransaction(tx *Transaction) bool {
	if tx.IsMinerTx() {
		return true
	}
	prevTXs := make(map[string]Transaction)

	for _, in := range tx.Inputs {
		prevTX, err := bc.FindTransaction(in.ID)
		Handle(err)
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}

	return tx.Verify(prevTXs)
}

func retry(dir string, originalOpts badger.Options) (*badger.DB, error) {
	lockPath := filepath.Join(dir, "LOCK")
	if err := os.Remove(lockPath); err != nil {
		return nil, fmt.Errorf(`removing "LOCK": %s`, err)
	}

	retryOpts := originalOpts
	retryOpts.Truncate = true
	db, err := badger.Open(retryOpts)
	return db, err
}

func OpenDB(dir string, opts badger.Options) (*badger.DB, error) {
	if db, err := badger.Open(opts); err != nil {
		if strings.Contains(err.Error(), "LOCK") {
			if db, err := retry(dir, opts); err == nil {
				log.Panicln("database unlocked , value log truncated ")
				return db, nil
			}

			log.Panicln("could not unlock databse", err)
		}

		return nil, err
	} else {
		return db, nil
	}
}
