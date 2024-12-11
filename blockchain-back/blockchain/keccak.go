package blockchain

import (
	"fmt"
	"log"
	"math"
	"math/big"
	"sync"

	"golang.org/x/crypto/sha3"
)

func (pow *ProofOfWork) KeccakLowRun() (int, []byte) {
	var intHash big.Int
	var hash []byte
	nonce := 0

	fmt.Println("\n-Low- Loading................")

	for nonce < math.MaxInt64 {
		data := pow.InitData(nonce)
		hash256 := sha3.NewLegacyKeccak256()

		if _, err := hash256.Write(data); err != nil {
			log.Panic("Error while hashing data: ", err)
		}

		hash = hash256.Sum(nil)
		intHash.SetBytes(hash)

		if intHash.Cmp(pow.Target) == -1 {
			break
		}

		nonce++
	}

	fmt.Println()
	return nonce, hash
}
func (pow *ProofOfWork) KeccakRun() (int, []byte) {

	numCPUs := 4

	fmt.Println("\n-High- Loading................")

	rangeSize := math.MaxInt64 / numCPUs
	var resultNonce int
	var resultHash []byte
	var wg sync.WaitGroup
	resultChan := make(chan struct {
		nonce int
		hash  []byte
	}, numCPUs)

	intHashPool := &sync.Pool{
		New: func() interface{} {
			return new(big.Int)
		},
	}
	hashPool := &sync.Pool{
		New: func() interface{} {
			return make([]byte, 32)
		},
	}

	for i := 0; i < numCPUs; i++ {
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()

			intHash := intHashPool.Get().(*big.Int)
			hash := hashPool.Get().([]byte)

			for nonce := start; nonce < end; nonce++ {
				data := pow.InitData(nonce)

				hash256 := sha3.NewLegacyKeccak256()
				if _, err := hash256.Write(data); err != nil {
					log.Panic("Error while hashing data: ", err)
				}

				hash = hash256.Sum(hash[:0])
				intHash.SetBytes(hash)

				if intHash.Cmp(pow.Target) == -1 {
					select {
					case resultChan <- struct {
						nonce int
						hash  []byte
					}{nonce: nonce, hash: hash[:]}:
						intHashPool.Put(intHash)
						hashPool.Put(hash)
						return
					}
				}
			}

			intHashPool.Put(intHash)
			hashPool.Put(hash)
		}(i*rangeSize, (i+1)*rangeSize)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	result := <-resultChan
	resultNonce = result.nonce
	resultHash = result.hash

	return resultNonce, resultHash
}

func (pow *ProofOfWork) KeccakValidate() bool {
	var intHash big.Int
	hash256 := sha3.NewLegacyKeccak256()

	data := pow.InitData(pow.Block.Nonce)

	if _, err := hash256.Write(data); err != nil {
		log.Panic("Error while hashing data: ", err)
	}

	hash := hash256.Sum(nil)

	intHash.SetBytes(hash)

	return intHash.Cmp(pow.Target) == -1
}
