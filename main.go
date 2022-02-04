package main

import (
	"flag"
	"fmt"
	"go-blockchain/blockchain"
	"os"
	"runtime"
	"strconv"
)

// Commandline struct for perform CLI operations
type Commandline struct {
	blockchain *blockchain.BlockChain
}

// cli options
func (cli *Commandline) printUsage() {
	fmt.Println("Usage:")
	fmt.Println(" add -block BLOCK_DATA -add a block to the chain.")
	fmt.Println(" print - Prints the blocks in the chain.")
}

// validate args passed
func (cli *Commandline) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		runtime.Goexit()
	}
}

// add block to chain
func (cli *Commandline) addBlock(data string) {
	cli.blockchain.AddBlock(data)
	fmt.Println("Added Block")
}

// print the chain
func (cli *Commandline) printChain() {
	iter := cli.blockchain.Iterator()

	for {
		block := iter.Next()

		fmt.Printf("Prev Hash: %x\n", block.PrevHash)
		fmt.Printf("Data in block: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)

		pow := blockchain.NewProof(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

	}
}

func (cli *Commandline) run() {
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
	// exit with status code 0 when main function is done
	defer os.Exit(0)
	chain := blockchain.InitBlockChain()

	// close the DB connection when main functions is done
	defer chain.Database.Close()

	// creat the cli
	cli := Commandline{chain}
	cli.run()

}