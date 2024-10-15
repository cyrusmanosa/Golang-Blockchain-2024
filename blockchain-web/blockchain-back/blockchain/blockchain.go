package blockchain

import (
	"fmt"
	"log"

	models "blockchain-back/modules"

	"github.com/dgraph-io/badger"
)

const (
	dbPath = "/Users/cyrusman/Desktop/ProgrammingLearning/GolangBlockchain2024/blockchain-web/blockchain-back/tmp/block"
)

type BlockChain struct {
	LastHash []byte
	Database *badger.DB
}

type BlockChainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

func (chain *BlockChain) Iterator() *BlockChainIterator {
	return &BlockChainIterator{chain.LastHash, chain.Database}
}

func (iter *BlockChainIterator) Next() *Block {
	var block *Block
	err := iter.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.CurrentHash)
		if err != nil {
			return fmt.Errorf("Iterator Next txn Get error: %w", err)
		}
		encodedBlock, err := item.Value()
		if err != nil {
			return fmt.Errorf("Iterator Next item value error: %w", err)
		}
		block = Deserialize(encodedBlock)
		return nil
	})

	if err != nil {
		log.Println("Next db View error: ", err)
		return nil
	}

	iter.CurrentHash = block.PrevHash
	return block
}

///------------------------------------------------ Guest ---------------------------------------------

func (chain *BlockChain) AddBlockForGuest(data models.InputData) error {
	var lastHash []byte

	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		if err != nil {
			return fmt.Errorf("AddBlockForGuest db Get error: %w", err)
		}
		lastHash, err = item.Value()
		return err
	})

	if err != nil {
		return fmt.Errorf("AddBlockForGuest db View error: %w", err)
	}

	newBlock := CreateBlockForGuest(data, lastHash)

	err = chain.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			return fmt.Errorf("AddBlockForGuest txn Set error: %w", err)
		}

		err = txn.Set([]byte("lh"), newBlock.Hash)
		if err != nil {
			return fmt.Errorf("AddBlockForGuest txn Set last hash error: %w", err)
		}
		chain.LastHash = newBlock.Hash
		return nil
	})

	return err
}

func InitBlockChainForGuest() (*BlockChain, error) {
	var lastHash []byte

	opts := badger.DefaultOptions
	opts.Dir = dbPath
	opts.ValueDir = dbPath

	db, err := badger.Open(opts)
	if err != nil {
		return nil, fmt.Errorf("InitBlockChainForGuest badger Open error: %w", err)
	}

	err = db.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound {
			GenesisForGuest := GenesisForGuest()

			err = txn.Set(GenesisForGuest.Hash, GenesisForGuest.Serialize())
			if err != nil {
				return fmt.Errorf("InitBlockChainForGuest txn set error: %w", err)
			}

			err = txn.Set([]byte("lh"), GenesisForGuest.Hash)
			lastHash = GenesisForGuest.Hash
			return err
		} else if err != nil {
			return fmt.Errorf("InitBlockChainForGuest txn Get error: %w", err)
		} else {
			item, err := txn.Get([]byte("lh"))
			if err != nil {
				return fmt.Errorf("InitBlockChainForGuest txn Get error: %w", err)
			}
			lastHash, err = item.Value()
			return err
		}
	})

	if err != nil {
		db.Close()
		return nil, fmt.Errorf("InitBlockChainForGuest db Update error: %w", err)
	}

	blockchain := BlockChain{lastHash, db}
	return &blockchain, nil
}

// /----------------------------------------- Doc ------------------------------------------------------
func InitBlockChainForDoc() *BlockChain {
	var lastHash []byte

	opts := badger.DefaultOptions
	opts.Dir = dbPath
	opts.ValueDir = dbPath

	db, err := badger.Open(opts)
	if err != nil {
		log.Println("InitBlockChainForDoc badger Open error: ", err)
	}

	err = db.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound {
			fmt.Println("No existing blockchain found")

			GenesisForGuest := GenesisForDoc()
			fmt.Println("GenesisForGuest proved")

			err = txn.Set(GenesisForGuest.Hash, GenesisForGuest.Serialize())
			if err != nil {
				log.Println("InitBlockChainForDoc txn set error: ", err)
			}

			err = txn.Set([]byte("lh"), GenesisForGuest.Hash)
			lastHash = GenesisForGuest.Hash
			return err
		} else {
			item, err := txn.Get([]byte("lh"))
			if err != nil {
				log.Println("InitBlockChainForDoc txn Get error: ", err)
			}

			lastHash, err = item.Value()
			return err
		}
	})

	if err != nil {
		log.Println("InitBlockChainForDoc db Update error: ", err)
	}

	blockchain := BlockChain{lastHash, db}
	return &blockchain
}

func (chain *BlockChain) AddBlockForDoc(data string) {
	var lastHash []byte

	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		if err != nil {
			log.Println("AddBlockForDoc txn Get error: ", err)
		}
		lastHash, err = item.Value()

		return err
	})
	if err != nil {
		log.Println("AddBlockForDoc db View error: ", err)
	}
	newBlock := CreateBlockForDoc(data, lastHash)

	err = chain.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			log.Println("AddBlockForDoc txn Set error: ", err)
		}
		err = txn.Set([]byte("lh"), newBlock.Hash)

		chain.LastHash = newBlock.Hash

		return err
	})

	if err != nil {
		log.Println("AddBlockForDoc db Update error: ", err)
	}
}
