package main

import (
	"os"

	"github.com/zhaizhonghao/blockchain-in-go/cli"
)

func main() {
	defer os.Exit(0)
	cli := cli.CommandLine{}
	cli.Run()

}
