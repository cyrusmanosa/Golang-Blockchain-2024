package blockchain

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/big"
	"runtime"
	"sync"

	"github.com/minio/highwayhash"
)

func (pow *ProofOfWork) HighWayHashLowRun() (int, []byte) {
	var intHash big.Int
	var hash []byte
	nonce := 0
	key := make([]byte, 32)

	fmt.Println("\n-Low- Loading................")

	hasher, err := highwayhash.New(key)
	if err != nil {
		log.Fatalf("Failed to create HighwayHash instance: %v", err)
	}

	for nonce < math.MaxInt64 {
		data := pow.InitData(nonce)
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

func (pow *ProofOfWork) HighWayHashRun() (int, []byte) {
	numCPUs := 4
	runtime.GOMAXPROCS(numCPUs)
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

			for nonce := start; nonce < end; nonce++ {
				select {
				case <-ctx.Done():
					return
				default:
					data := pow.InitData(nonce)
					hash := highwayhash.Sum(data, make([]byte, 32))
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

func (pow *ProofOfWork) HighWayHashValidate() bool {
	var intHash big.Int

	data := pow.InitData(pow.Block.Nonce)
	hash := highwayhash.Sum(data, make([]byte, 32))
	intHash.SetBytes(hash[:])

	return intHash.Cmp(pow.Target) == -1
}
