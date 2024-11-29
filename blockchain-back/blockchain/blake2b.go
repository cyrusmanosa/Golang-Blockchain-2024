package blockchain

import (
	"fmt"
	"log"
	"math"
	"math/big"
	"runtime"

	"golang.org/x/crypto/blake2b"
)

func (pow *ProofOfWork) Blake2bLowRun() (int, []byte) {
	var intHash big.Int
	var hash []byte
	nonce := 0

	runtime.GC()
	fmt.Println("\n-Low- Loading................")

	for nonce < math.MaxInt64 {
		data := pow.InitData(nonce)
		hasher, err := blake2b.New256(nil)
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
func (pow *ProofOfWork) Blake2bRun() (int, []byte) {
	numCPUs := 4

	var resultNonce int
	var resultHash []byte
	stopChan := make(chan struct{})
	resultChan := make(chan struct {
		nonce int
		hash  []byte
	})

	rangeSize := math.MaxInt32 / numCPUs
	runtime.GC()
	fmt.Println("\n-High- Loading................")

	for i := 0; i < numCPUs; i++ {
		go func(start, end int) {
			var intHash big.Int
			var hash []byte

			for nonce := start; nonce < end; nonce++ {
				select {
				case <-stopChan:
					return
				default:
					data := pow.InitData(nonce)

					hasher, err := blake2b.New256(nil)
					if err != nil {
						log.Panic("Failed to initialize Blake2b hasher: ", err)
					}

					if _, err := hasher.Write(data); err != nil {
						log.Panic("Error while hashing data: ", err)
					}

					hash = hasher.Sum(nil)
					intHash.SetBytes(hash[:])

					if intHash.Cmp(pow.Target) == -1 {
						select {
						case resultChan <- struct {
							nonce int
							hash  []byte
						}{nonce: nonce, hash: hash[:]}:
						case <-stopChan:
						}
						return
					}
				}
			}
		}(i*rangeSize, (i+1)*rangeSize)
	}

	result := <-resultChan
	resultNonce = result.nonce
	resultHash = result.hash
	close(stopChan)

	return resultNonce, resultHash
}

func (pow *ProofOfWork) Blake2bValidate() bool {
	var intHash big.Int
	hasher, err := blake2b.New256(nil)
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
