package blockchain

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"sync"

	"github.com/OneOfOne/xxhash"
)

func (pow *ProofOfWork) xxHash256(data []byte) []byte {
	hash := xxhash.New64()
	hash.Write(data)
	hashBytes := hash.Sum(nil)

	hash2 := xxhash.New64()
	hash2.Write(hashBytes)
	hash2Bytes := hash2.Sum(nil)

	hash3 := xxhash.New64()
	hash3.Write(hash2Bytes)
	hash3Bytes := hash3.Sum(nil)

	hash4 := xxhash.New64()
	hash4.Write(hash3Bytes)
	hash4Bytes := hash4.Sum(nil)

	result := append(hashBytes, hash2Bytes...)
	result = append(result, hash3Bytes...)
	result = append(result, hash4Bytes...)

	return result
}

func (pow *ProofOfWork) XxHashLowRun() (int, []byte) {
	var intHash big.Int
	var hash []byte

	nonce := 0

	fmt.Println("\n-Low- Loading................")

	for nonce < math.MaxInt64 {
		data := pow.InitData(nonce)
		hash = pow.xxHash256(data)
		intHash.SetBytes(hash[:])

		if intHash.Cmp(pow.Target) == -1 {
			break
		}
		nonce++
	}

	return nonce, hash[:]
}
func (pow *ProofOfWork) XxHashRun() (int, []byte) {
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
					hash := pow.xxHash256(data)

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

func (pow *ProofOfWork) XxhashValidate() bool {
	var intHash big.Int

	data := pow.InitData(pow.Block.Nonce)
	hash := pow.xxHash256(data)

	intHash.SetBytes(hash[:])

	return intHash.Cmp(pow.Target) == -1
}
