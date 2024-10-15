package controllers

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"blockchain-back/blockchain"
)

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println(" add -block BLOCK_DATA - add a block to the chain")
	fmt.Println(" print - Prints the blocks in the chain")
}

func validateArgs() {
	if len(os.Args) < 2 {
		printUsage()
		runtime.Goexit()
	}
}

func (cli *CommandLine) AddBlockForGuestForCommand(data string) {
	cli.blockchain.AddBlockForDoc(data)
	fmt.Println("Added Block!")
}

func (cli *CommandLine) CommandRun() {
	for {
		fmt.Println("Please choose a function to run")
		fmt.Println("1. Add  2. Print")
		validateArgs()

		AddBlockForGuestCmd := flag.NewFlagSet("add", flag.ExitOnError)
		printChainCmd := flag.NewFlagSet("print", flag.ExitOnError)
		AddBlockForGuestData := AddBlockForGuestCmd.String("block", "", "Block data")

		var choice string
		fmt.Scan(&choice)

		switch choice {
		case "1":
			// add date
			var blockData string
			fmt.Print("Enter block data: ")
			fmt.Scan(&blockData)

			os.Args = append([]string{os.Args[0], "add", "-block=" + blockData}, os.Args[2:]...)
			err := AddBlockForGuestCmd.Parse(os.Args[2:])
			if err != nil {
				fmt.Println("Error parsing 'add' command arguments:", err)
				AddBlockForGuestCmd.Usage()
				runtime.Goexit()
			}

			if *AddBlockForGuestData == "" {
				fmt.Println("No block data provided. Use -block flag to specify block data.")
				AddBlockForGuestCmd.Usage()
				runtime.Goexit()
			}
			cli.AddBlockForGuestForCommand(*AddBlockForGuestData)
		case "2":
			err := printChainCmd.Parse(os.Args[2:])
			if err != nil {
				log.Println("2 error: ", err)
			}
			if printChainCmd.Parsed() {
				cli.PrintChain()
			}

		default:
			fmt.Println("Invalid option. Please choose again.")
		}
	}
}

func Cli() {
	chain := blockchain.InitBlockChainForDoc()
	defer chain.Database.Close()

	cli := CommandLine{chain}
	cli.CommandRun()
}
