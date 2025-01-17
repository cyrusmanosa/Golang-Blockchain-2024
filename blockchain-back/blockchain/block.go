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

// /------------------------------------------------ Guest ------------------------------------------------------
func CreateBlockForGuest(data models.InputData, prevHash []byte) *Block {
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshalling to JSON:", err)
	}

	block := &Block{[]byte{}, []byte(jsonData), prevHash, 0}
	pow := NewProof(block)
	///------ ************************************ ------
	// if len(data.File) > 1024*1024 {
	switch data.Hash {
	case "sha256", "":
		nonce, hash := pow.Sha256Run()
		// nonce, hash := pow.Sha256LowRun()
		block.result(nonce, hash)
	case "blake2b":
		nonce, hash := pow.Blake2bRun()
		block.result(nonce, hash)
	case "blake3":
		nonce, hash := pow.Blake3Run()
		block.result(nonce, hash)
	case "murmurHash3":
		// nonce, hash := pow.MurmurHashRun()
		nonce, hash := pow.MurmurHashLowRun()
		block.result(nonce, hash)
	case "keccak":
		nonce, hash := pow.KeccakRun()
		block.result(nonce, hash)
	case "skein":
		nonce, hash := pow.SkeinRun()
		block.result(nonce, hash)
	case "farmHash":
		nonce, hash := pow.FarmRun()
		block.result(nonce, hash)
	case "xxHash":
		nonce, hash := pow.XxHashRun()
		block.result(nonce, hash)
	case "highwayHash":
		nonce, hash := pow.HighWayHashRun()
		block.result(nonce, hash)
	}
	return block
}

func (b *Block) result(nonce int, hash []byte) *Block {
	b.Hash = hash[:]
	b.Nonce = nonce
	return b
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
		// File:        "",
		File:     nil,
		Status:   "New One",
		SendTime: time,
	}

	return CreateBlockForGuest(GenesisForGuest, []byte{})
}

// asdasdas
