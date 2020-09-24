package blockchain

import (
	"bytes"
	"encoding/gob"
)

//must ensure that at least one transaction is stored in one block
type Block struct {
	Hash         []byte
	Transactions []*Transaction
	PrevHash     []byte
	Nonce        int
}

//To create the gensis block
func Genesis(coinbase *Transaction) *Block {
	return CreateBlock([]*Transaction{coinbase}, []byte{})
}

func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.Serialize())
	}
	tree := NewMerkleTree(txHashes)

	return tree.RootNode.Data
}

func CreateBlock(txs []*Transaction, prevHash []byte) *Block {
	block := &Block{[]byte{}, txs, prevHash, 0}
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
