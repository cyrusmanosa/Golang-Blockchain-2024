package blockchain

import (
	models "Sha256v/modules"
	"encoding/json"
	"fmt"
	"time"
)

type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
	Nonce    int
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
	nonce, hash := pow.Run()
	///------ ************************************ ------

	block.Hash = hash[:]
	block.Nonce = nonce

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
		// File:        nil,
		File:        "",
		Status:      "New One",
		SendTime:    time,
		ConfirmTime: "",
	}

	return CreateBlockForGuest(GenesisForGuest, []byte{})
}

///------------------------------------------------ Doc ------------------------------------------------------

func CreateBlockForDoc(data string, prevHash []byte) *Block {
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshalling to JSON:", err)
	}

	block := &Block{[]byte{}, []byte(jsonData), prevHash, 0}
	pow := NewProof(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

func GenesisForDoc() *Block {
	GenesisForDoc := "GenesisForGuest"
	return CreateBlockForDoc(GenesisForDoc, []byte{})
}
