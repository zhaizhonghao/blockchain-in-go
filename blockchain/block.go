package blockchain

import (
	"bytes"
	"encoding/gob"
	"log"

	"github.com/dgraph-io/badger"
)

type Blockchain struct {
	LastHash []byte
	Database *badger.DB
}

type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
	Nonce    int
}

func InitBlockchain() *Blockchain {
	return &Blockchain{[]*Block{Genesis()}}
}

//To create the gensis block
func Genesis() *Block {
	return CreateBlock("Genesis", []byte{})
}

func (chain *Blockchain) AddBlock(data string) {
	prevBlock := chain.Blocks[len(chain.Blocks)-1]
	new := CreateBlock(data, prevBlock.Hash)
	chain.Blocks = append(chain.Blocks, new)
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

//serialize the block
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
