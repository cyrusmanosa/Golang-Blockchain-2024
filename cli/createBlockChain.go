package cli

import (
	"fmt"
	"golang-blockchain/blockchain"
	"golang-blockchain/wallet"
	"log"
)

func (cli *CommandLine) CreateBlockChain(address, nodeID string) {
	if !wallet.ValidateAddress(address) {
		log.Panic("Address is not Valid")
	}
	chain := blockchain.InitBlockChain(address, nodeID)
	defer chain.Database.Close()

	UTXOSet := blockchain.UTXOSet{Blockchain: chain}
	UTXOSet.Reindex()

	fmt.Println("Finished!")
}
