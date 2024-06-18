package cli

import (
	"fmt"
	"golang-blockchain/wallet"
)

func (cli *CommandLine) ListAddresses(nodeID string) {
	wallets, _ := wallet.CreateWallets(nodeID)
	addresses := wallets.GetAllAddresses()

	for _, address := range addresses {
		fmt.Println(address)
	}
}
