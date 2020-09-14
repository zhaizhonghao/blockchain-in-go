package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
)

type Blockchain struct {
	blocks []*Block
}

type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
}

func (b *Block) DeriveHash() {
	//func Join(s [][]byte, sep []byte) []byte
	//Join concatenates the elements of s to create a new byte slice. The separator sep is placed between elements in the resulting slice. Here the separator is {}
	info := bytes.Join([][]byte{b.Data, b.PrevHash}, []byte{})
	//compute the hash of the block with SHA-256
	hash := sha256.Sum256(info)
	b.Hash = hash[:]
}

func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{[]byte{}, []byte(data), prevHash}
	block.DeriveHash()
	return block
}

func (chain *Blockchain) AddBlock(data string) {
	prevBlock := chain.blocks[len(chain.blocks)-1]
	new := CreateBlock(data, prevBlock.Hash)
	chain.blocks = append(chain.blocks, new)
}

//To create the gensis block
func Genesis() *Block {
	return CreateBlock("Genesis", []byte{})
}

func InitBlockchain() *Blockchain {
	return &Blockchain{[]*Block{Genesis()}}
}

func main() {
	fmt.Println("Hello blockchain!")
	chain := InitBlockchain()

	chain.AddBlock("First Block after Genesis")
	chain.AddBlock("Second Block")
	chain.AddBlock("Third Block")

	for _, block := range chain.blocks {
		fmt.Printf("Previous Hash:%x\n", block.PrevHash)
		fmt.Printf("Data in Block:%s\n", block.Data)
		fmt.Printf("Hash:%x\n", block.Hash)
		fmt.Println()
	}
}
