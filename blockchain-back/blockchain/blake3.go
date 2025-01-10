package blockchain

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/big"
	"runtime"
	"sync"

	"lukechampine.com/blake3"
)

func (pow *ProofOfWork) Blake3LowRun() (int, []byte) {
	var intHash big.Int
	var hash []byte
	nonce := 0

	fmt.Println("\n-Low- Loading................")

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
func (pow *ProofOfWork) Blake3Run() (int, []byte) {
	numCPUs := 4
	runtime.GOMAXPROCS(numCPUs)
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

					hasher := blake3.New(32, nil)
					if _, err := hasher.Write(data); err != nil {
						log.Panic("Error while hashing data: ", err)
					}

					hash := hasher.Sum(nil)
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
