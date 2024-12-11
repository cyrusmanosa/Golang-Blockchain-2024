package blockchain

import (
	"context"
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

	var once sync.Once
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	resultChan := make(chan struct {
		nonce int
		hash  []byte
	}, numCPUs)

	rangeSize := math.MaxInt64 / numCPUs
	wg.Add(numCPUs)

	for i := 0; i < numCPUs; i++ {
		start := i * rangeSize
		end := start + rangeSize

		go func(start, end int) {
			defer wg.Done()
			var intHash big.Int

			for nonce := start; nonce < end; nonce++ {
				select {
				case <-ctx.Done():
					return
				default:
					data := pow.InitData(nonce)
					hash256 := sha3.NewLegacyKeccak256()
					if _, err := hash256.Write(data); err != nil {
						log.Panic("Error while hashing data: ", err)
					}
					hash := hash256.Sum(nil)
					intHash.SetBytes(hash)

					if intHash.Cmp(pow.Target) == -1 {
						once.Do(func() {
							resultChan <- struct {
								nonce int
								hash  []byte
							}{nonce: nonce, hash: hash}
							cancel()
						})
						return
					}
				}
			}
		}(start, end)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	result := <-resultChan
	return result.nonce, result.hash
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
