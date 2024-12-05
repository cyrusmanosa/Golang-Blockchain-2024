package blockchain

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/big"
	"runtime"
	"sync"

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
	runtime.GC()
	numCPUs := 4

	fmt.Println("\n-High- Loading................")

	var resultNonce int
	var resultHash []byte
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
			var hash []byte

			for nonce := start; nonce < end; nonce++ {
				select {
				case <-ctx.Done():
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
						once.Do(func() {
							resultChan <- struct {
								nonce int
								hash  []byte
							}{nonce: nonce, hash: hash[:]}
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

	result, ok := <-resultChan
	if ok {
		resultNonce = result.nonce
		resultHash = result.hash
	}

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
