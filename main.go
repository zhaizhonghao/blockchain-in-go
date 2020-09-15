package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"

	"github.com/zhaizhonghao/blockchain-in-go/blockchain"
)

type CommandLine struct {
	blockchain *blockchain.Blockchain
}

func (cli *CommandLine) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("add -block BLOCK_DATA - add a block to the chain")
	fmt.Println("print - Prints the blocks in the chain")
}

func (cli *CommandLine) validateArgs() {
	//Args hold the command-line arguments, starting with the program name.
	if len(os.Args) < 2 {
		cli.printUsage()
		//exit all goroutines in the application, especailly the database to protect the data
		runtime.Goexit()
	}
}
func (cli *CommandLine) addBlock(data string) {
	cli.blockchain.AddBlock(data)
	fmt.Println("adding a block!")
}

func (cli *CommandLine) printChain() {
	iter := cli.blockchain.Iterator()

	for {
		block := iter.Next()

		fmt.Printf("Previous Hash:%x\n", block.PrevHash)
		fmt.Printf("Data in Block:%s\n", block.Data)
		fmt.Printf("Hash:%x\n", block.Hash)

		pow := blockchain.NewProof(block)
		fmt.Printf("Valid PoW:%s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevHash) == 0 {
			break
		}
	}
}

func (cli *CommandLine) run() {
	cli.validateArgs()

	addBlockCmd := flag.NewFlagSet("add", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("print", flag.ExitOnError)
	addBlockData := addBlockCmd.String("block", "", "Block data")

	switch os.Args[1] {
	case "add":
		err := addBlockCmd.Parse(os.Args[2:])
		blockchain.Handle(err)

	case "print":
		err := printChainCmd.Parse(os.Args[2:])
		blockchain.Handle(err)

	default:
		cli.printUsage()
		runtime.Goexit()
	}

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			runtime.Goexit()
		}
		cli.addBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

}

func main() {
	defer os.Exit(0)
	chain := blockchain.InitBlockchain()
	defer chain.Database.Close()

	cli := CommandLine{chain}
	cli.run()
}
