package controllers

import (
	blockchain "blockchain-back/blockchain"
	"flag"
	"fmt"
	"log"
	"os"
)

func (cli *CommandLine) AddBlockForGuestForCommand(data string) {
	cli.blockchain.AddBlockForDoc(data)
	fmt.Println("Added Block!")
}

func (cli *CommandLine) CommandRun() {

	printChainCmd := flag.NewFlagSet("print", flag.ExitOnError)
	err := printChainCmd.Parse(os.Args[2:])
	if err != nil {
		log.Println("2 error: ", err)
	}
	if printChainCmd.Parsed() {
		cli.PrintChain()
	}
}

func Cli() {
	chain := blockchain.InitBlockChainForDoc()
	defer chain.Database.Close()

	cli := CommandLine{chain}
	cli.CommandRun()
}
