package main

import (
	"fmt"
	"strconv"

	"github.com/zhaizhonghao/blockchain-in-go/blockchain"
)

func main() {
	fmt.Println("Hello blockchain!")
	chain := blockchain.InitBlockchain()

	chain.AddBlock("First Block after Genesis")
	chain.AddBlock("Second Block")
	chain.AddBlock("Third Block")

	for _, block := range chain.Blocks {
		fmt.Printf("Previous Hash:%x\n", block.PrevHash)
		fmt.Printf("Data in Block:%s\n", block.Data)
		fmt.Printf("Hash:%x\n", block.Hash)

		pow := blockchain.NewProof(block)
		fmt.Printf("Valid PoW:%s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
	}
}
