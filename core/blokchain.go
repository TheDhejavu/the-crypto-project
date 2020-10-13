package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	badger "github.com/dgraph-io/badger"
	log "github.com/sirupsen/logrus"
)

// Blockchain struct such that lastHash represents the lastblock hash
// on the ledger
type Blockchain struct {
	LastHash   []byte
	Database   *badger.DB
	InstanceId string
}

var (
	mutex      = &sync.Mutex{}
	_, b, _, _ = runtime.Caller(0)

	// Root folder of this project
	Root        = filepath.Join(filepath.Dir(b), "../")
	genesisData = "genesis"
)

// Check if Blockchain Database already exist
func DBExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func Exists(instanceId string) bool {
	return DBExists(GetDatabasePath(instanceId))
}

func GetDatabasePath(port string) string {
	if port != "" {
		return filepath.Join(Root, fmt.Sprintf("./tmp/blocks_%s", port))
	}
	return filepath.Join(Root, "./tmp/blocks")
}

func OpenBardgerDB(instanceId string) *badger.DB {
	path := GetDatabasePath(instanceId)

	// if DBExists(path) == false {
	// 	log.Info("No Existing Blockchian DB found, create one!")
	// 	runtime.Goexit()
	// }

	opts := badger.DefaultOptions(path)
	db, err := OpenDB(path, opts)
	Handle(err)

	return db
}

func (chain *Blockchain) ContinueBlockchain() *Blockchain {
	var lastHash []byte
	var db *badger.DB
	if chain.Database == nil {
		db = OpenBardgerDB(chain.InstanceId)
	} else {
		db = chain.Database
	}

	//Read-Write Operations
	err := db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		if err == nil {
			lastHash, err = item.ValueCopy(nil)
		}

		return err
	})

	if err != nil {
		lastHash = nil
	}
	// log.Infof("LastHash: %x", lastHash)
	return &Blockchain{lastHash, db, chain.InstanceId}
}

// Initialize the blockchain by creating the blockchain database
// with a genesis block with an address
func InitBlockchain(address string, instanceId string) *Blockchain {
	var lastHash []byte
	path := GetDatabasePath(instanceId)

	if DBExists(path) {
		log.Info("Blockchain already exist")
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
		log.Info("No existing blockchain found")
		genesis := Genesis(cbtx)
		err = txn.Set(genesis.Hash, genesis.Serialize())
		Handle(err)
		err = txn.Set([]byte("lh"), genesis.Hash)
		lastHash = genesis.Hash

		return err
	})
	Handle(err)

	return &Blockchain{lastHash, db, instanceId}
}

// Add a block to the blockchain
//https://github.com/dgraph-io/badger#read-write-transactions
func (chain *Blockchain) AddBlock(block *Block) *Block {
	mutex.Lock()

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
		if err == nil {
			lastHash, _ := item.ValueCopy(nil)
			item, err = txn.Get(lastHash)
			Handle(err)
			lastBlockData, _ := item.ValueCopy(nil)
			lastBlock := DeSerialize(lastBlockData)

			// check if the current block height is
			// greater than the lastBlock Height
			if block.Height > lastBlock.Height {
				err := txn.Set([]byte("lh"), block.Hash)
				Handle(err)
				chain.LastHash = block.Hash
			}
		} else {
			err = txn.Set([]byte("lh"), block.Hash)
			chain.LastHash = block.Hash
		}

		return err
	})

	Handle(err)
	mutex.Unlock()
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
func (chain *Blockchain) GetBlockHashes(height int) [][]byte {
	var blocks [][]byte

	iter := chain.Iterator()
	if iter == nil {
		return blocks
	}
	for {
		block := iter.Next()
		prevHash := block.PrevHash
		if block.Height == height {
			break
		}
		blocks = append([][]byte{block.Hash}, blocks...)

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
		if err == nil {
			lastHash, _ := item.ValueCopy(nil)

			item, err = txn.Get(lastHash)
			Handle(err)
			lastBlockData, _ := item.ValueCopy(nil)
			lastBlock = *DeSerialize(lastBlockData)
		}

		return err
	})

	if err == nil {
		return lastBlock.Height
	}

	return 0
}

//Mine Block Creates a new block and add it to the blockchain
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

//Find a specific transaction by ID
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
	log.Error("Error: No Transaction with ID")

	return Transaction{}, errors.New("No transaction with id")
}

func (chain *Blockchain) GetTransaction(transaction *Transaction) map[string]Transaction {
	txs := make(map[string]Transaction)
	for _, in := range transaction.Inputs {
		// get all transaction with in.ID
		tx, err := chain.FindTransaction(in.ID)
		if err != nil {
			log.Error("Error: Invalid Transaction Ewwww")
		}
		Handle(err)
		txs[hex.EncodeToString(tx.ID)] = tx
	}

	return txs
}

func (chain *Blockchain) SignTransaction(privKey ecdsa.PrivateKey, tx *Transaction) {
	prevTxs := chain.GetTransaction(tx)
	tx.Sign(privKey, prevTxs)
}

func (chain *Blockchain) VerifyTransaction(tx *Transaction) bool {
	if tx.IsMinerTx() {
		return true
	}
	prevTxs := chain.GetTransaction(tx)

	return tx.Verify(prevTxs)
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

			log.Panicln("could not unlock database", err)
		}

		return nil, err
	} else {
		return db, nil
	}
}
