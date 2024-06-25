package cli

import (
	"GolangBlockchain2024/blockchain"
	"GolangBlockchain2024/wallet"
	"fmt"
	"log"
)

func (cli *CommandLine) createBlockChain(address, nodeID string) {
	if !wallet.ValidateAddress(address) {
		log.Panic("Address is not Valid")
	}
	chain := blockchain.InitBlockChain(address, nodeID)
	defer chain.Database.Close()

	UTXOSet := blockchain.UTXOSet{Blockchain: chain}
	UTXOSet.Reindex()

	fmt.Println("Finished!")
}
