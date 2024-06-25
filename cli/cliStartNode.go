package cli

import (
	"GolangBlockchain2024/network"
	"GolangBlockchain2024/wallet"
	"fmt"
	"log"
)

func (cli *CommandLine) StarttNode(nodeID, minerAddress string) {
	fmt.Printf("Starting Node %s\n", nodeID)

	if len(minerAddress) > 0 {
		if wallet.ValidateAddress(minerAddress) {
			fmt.Println("Mining is on. Address to receive rewards: ", minerAddress)
		} else {
			log.Panic("Wrong miner address!")
		}
	}
	network.StartServer(nodeID, minerAddress)
}
