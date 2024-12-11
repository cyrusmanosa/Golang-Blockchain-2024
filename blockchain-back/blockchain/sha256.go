package blockchain

import (
	"context"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
	"sync"
)

func (pow *ProofOfWork) Sha256LowRun() (int, []byte) {
	var intHash big.Int
	var hash [32]byte

	nonce := 0

	fmt.Println("\n-Low- Loading................")
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

	return nonce, hash[:]
}

func (pow *ProofOfWork) Sha256Run() (int, []byte) {

	numCPUs := 4

	fmt.Println("\n-High- Loading................")

	var resultNonce int
	var resultHash []byte
	var once sync.Once
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rangeSize := math.MaxInt64 / numCPUs
	resultChan := make(chan struct {
		nonce int
		hash  []byte
	}, numCPUs)

	wg.Add(numCPUs)
	for i := 0; i < numCPUs; i++ {
		start := i * rangeSize
		end := start + rangeSize
		go func(start, end int) {
			defer wg.Done()

			var intHash big.Int
			var hash [32]byte

			for nonce := start; nonce < end; nonce++ {
				select {
				case <-ctx.Done():
					return
				default:
					data := pow.InitData(nonce)
					hash = sha256.Sum256(data)
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

func (pow *ProofOfWork) Sha256Validate() bool {
	var intHash big.Int

	data := pow.InitData(pow.Block.Nonce)
	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:])

	return intHash.Cmp(pow.Target) == -1
}
