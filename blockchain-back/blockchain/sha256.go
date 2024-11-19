package blockchain

import (
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)

func (pow *ProofOfWork) Sha256Run() (int, []byte) {
	var intHash big.Int
	var hash [32]byte

	nonce := 0
	fmt.Println("Loading................")
	for nonce < math.MaxInt64 {
		data := pow.InitData(nonce)
		hash = sha256.Sum256(data)
		// fmt.Printf("\r%x", hash)
		intHash.SetBytes(hash[:])

		if intHash.Cmp(pow.Target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Println()

	return nonce, hash[:]
}

func (pow *ProofOfWork) Sha256Validate() bool {
	var intHash big.Int

	data := pow.InitData(pow.Block.Nonce)

	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:])

	return intHash.Cmp(pow.Target) == -1
}
