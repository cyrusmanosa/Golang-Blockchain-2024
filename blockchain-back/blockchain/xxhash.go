package blockchain

import (
	"context"
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
	"runtime"
	"sync"

	"github.com/cespare/xxhash/v2"
)

func (pow *ProofOfWork) XxHashLowRun() (int, []byte) {
	var intHash big.Int
	var hash [32]byte

	nonce := 0
	runtime.GC()
	fmt.Println("\n-Low- Loading................")

	for nonce < math.MaxInt64 {
		data := pow.InitData(nonce)
		hash = xxHash256(data)
		intHash.SetBytes(hash[:])

		if intHash.Cmp(pow.Target) == -1 {
			break
		}
		nonce++
	}

	return nonce, hash[:]
}

func (pow *ProofOfWork) XxHashRun() (int, []byte) {
	runtime.GC()
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
					hash = xxHash256(data)
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

func (pow *ProofOfWork) XxhashValidate() bool {
	var intHash big.Int

	data := pow.InitData(pow.Block.Nonce)
	hash := xxHash256(data)

	intHash.SetBytes(hash[:])

	return intHash.Cmp(pow.Target) == -1
}

func xxHash256(data []byte) [32]byte {
	hash1 := xxhash.Sum64(data[:len(data)/2])
	hash2 := xxhash.Sum64(data[len(data)/2:])

	var hash256 [32]byte
	binary.LittleEndian.PutUint64(hash256[:8], hash1)
	binary.LittleEndian.PutUint64(hash256[8:16], hash2)
	binary.LittleEndian.PutUint64(hash256[16:24], hash1)
	binary.LittleEndian.PutUint64(hash256[24:], hash2)

	return hash256
}
