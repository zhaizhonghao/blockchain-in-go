package blockchain

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"

	"github.com/dgraph-io/badger"
)

const (
	//the path to store the databases
	dbPath = "./tmp/blocks"
)

type Blockchain struct {
	//Blocks []*Block

	//use the Database to store the blockchian
	LastHash []byte
	//a key-value database in golang
	Database *badger.DB
}

//To travase the blockchain in the database
type BlockchainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
	Nonce    int
}

func InitBlockchain() *Blockchain {
	//return &Blockchain{[]*Block{Genesis()}}

	var lastHash []byte

	opts := badger.DefaultOptions(dbPath)
	opts.Dir = dbPath
	opts.ValueDir = dbPath

	db, err := badger.Open(opts)
	Handle(err)

	//To write the database
	err = db.Update(func(txn *badger.Txn) error {
		//At first, try to fetch the last hash((lh) from the databased
		if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound {
			//there is no blockchain in our database
			fmt.Println("No existing blockchain found")
			genesis := Genesis()
			fmt.Println("Genesis block generated!")
			//store the genesis block in key-value
			//where the key is the hash of the block, and the value is the genesis block in bytes by serailizing
			err = txn.Set(genesis.Hash, genesis.Serialize())
			Handle(err)
			//set the last hash of the blockchain
			err = txn.Set([]byte("lh"), genesis.Hash)

			lastHash = genesis.Hash
			return err
		} else {
			//we already have the blockchain in our system
			//we can set the lastHash to the lh in the database
			item, err := txn.Get([]byte("lh"))
			Handle(err)
			err = item.Value(func(val []byte) error {
				lastHash = append([]byte{}, val...)
				return nil
			})
			return err
		}
	})

	Handle(err)

	blockchain := Blockchain{lastHash, db}

	return &blockchain

}

//To create the gensis block
func Genesis() *Block {
	return CreateBlock("Genesis", []byte{})
}

func (chain *Blockchain) AddBlock(data string) {
	var lastHash []byte
	//read the databased to fetch the last hash of the blockchain
	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		err = item.Value(func(val []byte) error {
			lastHash = append([]byte{}, val...)
			return nil
		})
		return err
	})
	Handle(err)

	newBlock := CreateBlock(data, lastHash)
	//append a block to the blockchain in the database, and set the lh to the hash of the new block
	err = chain.Database.Update(func(txn *badger.Txn) error {
		//append the new block in the database
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		Handle(err)
		//To set the lh to the hash of the new block
		err = txn.Set([]byte("lh"), newBlock.Hash)
		//change the LastHash of the blockchain in Memory
		chain.LastHash = newBlock.Hash

		return err
	})
}

func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{[]byte{}, []byte(data), prevHash, 0}
	pow := NewProof(block)
	//mining
	nonce, hash := pow.Run()
	//mining successfully
	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}

//serialize the block to the bytes
func (b *Block) Serialize() []byte {
	var res bytes.Buffer
	//An Encoder manages the transmission of type and data information to the other side of a connection.
	//It is safe for concurrent use by multiple goroutines.
	encoder := gob.NewEncoder(&res)
	//translate the block into bytes in gob
	err := encoder.Encode(b)
	Handle(err)

	return res.Bytes()
}

//deserialize the bytes to the block
func Deserialize(data []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&block)

	Handle(err)

	return &block
}

func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func (chain *Blockchain) Iterator() *BlockchainIterator {
	iter := &BlockchainIterator{chain.LastHash, chain.Database}
	return iter
}

//iterate the blockchain until the genesis block
func (iter *BlockchainIterator) Next() *Block {
	var block *Block
	err := iter.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.CurrentHash)
		Handle(err)
		err = item.Value(func(val []byte) error {
			encodeBlock := append([]byte{}, val...)
			block = Deserialize(encodeBlock)
			return nil
		})
		return err
	})
	Handle(err)

	iter.CurrentHash = block.PrevHash

	return block
}
