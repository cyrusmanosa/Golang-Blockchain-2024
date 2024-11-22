package blockchain

import (
	"encoding/json"
	"fmt"
	"time"

	models "blockchain-back/modules"
)

type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
	Nonce    int
}

///------------------------------------------------ Doc ------------------------------------------------------

func CreateBlockForDoc(data string, prevHash []byte) *Block {
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshalling to JSON:", err)
	}

	block := &Block{[]byte{}, []byte(jsonData), prevHash, 0}
	pow := NewProof(block)
	nonce, hash := pow.Sha256Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

func GenesisForDoc() *Block {
	GenesisForDoc := "GenesisForGuest"
	return CreateBlockForDoc(GenesisForDoc, []byte{})
}

///------------------------------------------------ Guest ------------------------------------------------------

func CreateBlockForGuest(data models.InputData, prevHash []byte) *Block {
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshalling to JSON:", err)
	}

	block := &Block{[]byte{}, []byte(jsonData), prevHash, 0}
	pow := NewProof(block)
	///------ ************************************ ------
	switch data.Hash {
	case "sha256", "":
		nonce, hash := pow.Sha256Run()
		block.Hash = hash[:]
		block.Nonce = nonce
	case "argon2":
		nonce, hash := pow.Argon2Run()
		block.Hash = hash[:]
		block.Nonce = nonce
	case "blake2b":
		nonce, hash := pow.Blake2bRun()
		block.Hash = hash[:]
		block.Nonce = nonce
		// case "blake2s":
		// 	nonce, hash := pow.Blake2sRun()
		// 	block.Hash = hash[:]
		// 	block.Nonce = nonce
	}
	///------ ************************************ ------
	return block
}

func GenesisForGuest() *Block {
	layout := "2006-01-02 15:04:05"
	time := time.Now().Format(layout)

	GenesisForGuest := models.InputData{
		Name:        "GenesisForGuest",
		Email:       "",
		CompanyName: "",
		Message:     "",
		Hash:        "",
		File:        "",
		Status:      "New One",
		SendTime:    time,
	}

	return CreateBlockForGuest(GenesisForGuest, []byte{})
}
