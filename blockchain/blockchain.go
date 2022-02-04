package blockchain

import (
	"fmt"
	"github.com/dgraph-io/badger"
)

// path to DB
const (
	dbPath = "./tmp/blocks"
)

// BlockChain is a basic structure for a blockchain
type BlockChain struct {
	LashHash []byte
	Database *badger.DB
}

func InitBlockChain() *BlockChain {
	var lastHash []byte

	opts := badger.DefaultOptions
	opts.Dir = dbPath
	opts.ValueDir = dbPath

	// open db connection
	db, err := badger.Open(opts)
	Handle(err)

	// passing in closure
	// check if there is a blockchain, if not create it
	err = db.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound {
			fmt.Println("No existing blockchain found")
			genesis := Genesis()
			fmt.Println("Genesis proved")
			err = txn.Set(genesis.Hash, genesis.Serialize())
			Handle(err)
			err = txn.Set([]byte("lh"), genesis.Hash)

			lastHash = genesis.Hash

			return err
		} else {
			item, err := txn.Get([]byte("lh"))
			Handle(err)
			lastHash, err = item.Value()
			return err
		}
	})

	// handle any error from the if/else
	Handle(err)

	blockchain := BlockChain{lastHash, db}
	return &blockchain
}

// AddBlock adds a block to the blockchain
func (chain *BlockChain) AddBlock(data string) {
	var lastHash []byte

	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		lastHash, err = item.Value()

		return err
	})
	Handle(err)

	newBlock := CreateBlock(data, lastHash)

	err = chain.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		Handle(err)
		err = txn.Set([]byte("lh"), newBlock.Hash)

		chain.LastHash = newBlock.Hash

		return err
	})
	Handle(err)
}

