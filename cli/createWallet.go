package cli

import (
	"fmt"
	"golang-blockchain/wallet"
)

func (cli *CommandLine) CreateWallet(nodeID string) {
	wallets, _ := wallet.CreateWallets(nodeID)
	address := wallets.AddWallet()
	wallets.SaveFile(nodeID)

	fmt.Printf("New address is: %s\n", address)
}
