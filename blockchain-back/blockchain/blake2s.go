package blockchain

import (
	"fmt"
	"log"
	"math"
	"math/big"
	"runtime"

	"golang.org/x/crypto/blake2s"
)

func (pow *ProofOfWork) Blake2sRun() (int, []byte) {
	var intHash big.Int
	var hash []byte
	nonce := 0

	runtime.GC()
	fmt.Println("Loading................")

	for nonce < math.MaxInt64 {
		data := pow.InitData(nonce)
		hasher, err := blake2s.New256(nil)
		if err != nil {
			log.Panic("Failed to initialize Blake2b hasher: ", err)
		}

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

func (pow *ProofOfWork) Blake2sValidate() bool {
	var intHash big.Int
	hasher, err := blake2s.New256(nil)
	if err != nil {
		log.Panic(err)
	}

	data := pow.InitData(pow.Block.Nonce)

	if _, err = hasher.Write(data); err != nil {
		log.Panic("Error while hashing data: ", err)
	}

	hash := hasher.Sum(nil)

	intHash.SetBytes(hash)

	return intHash.Cmp(pow.Target) == -1
}
