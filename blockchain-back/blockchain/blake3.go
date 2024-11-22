package blockchain

import (
	"fmt"
	"log"
	"math"
	"math/big"

	"lukechampine.com/blake3"
)

func (pow *ProofOfWork) Blake3Run() (int, []byte) {
	var intHash big.Int
	var hash []byte
	nonce := 0

	fmt.Println("Loading................")

	for nonce < math.MaxInt64 {
		data := pow.InitData(nonce)
		hasher := blake3.New(32, nil)

		if _, err := hasher.Write(data); err != nil {
			log.Panic("Error while hashing data: ", err)
		}

		hash = hasher.Sum(nil)
		intHash.SetBytes(hash)

		if intHash.Cmp(pow.Target) == -1 {
			break
		}

		nonce++
	}

	fmt.Println()
	return nonce, hash
}

func (pow *ProofOfWork) Blake3Validate() bool {
	var intHash big.Int
	hasher := blake3.New(32, nil)

	data := pow.InitData(pow.Block.Nonce)

	if _, err := hasher.Write(data); err != nil {
		log.Panic("Error while hashing data: ", err)
	}

	hash := hasher.Sum(nil)

	intHash.SetBytes(hash)

	return intHash.Cmp(pow.Target) == -1
}
