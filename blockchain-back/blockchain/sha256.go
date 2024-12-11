package blockchain

import (
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
	"runtime"
	"sync"
)

func (pow *ProofOfWork) Sha256LowRun() (int, []byte) {
	var intHash big.Int
	var hash [32]byte

	nonce := 0
	runtime.GC()
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

	intPool := &sync.Pool{
		New: func() interface{} {
			return &big.Int{}
		},
	}
	hashPool := &sync.Pool{
		New: func() interface{} {
			var hash [32]byte
			return &hash
		},
	}

	for i := 0; i < numCPUs; i++ {
		go func(start, end int) {
			intHash := intPool.Get().(*big.Int)
			hash := hashPool.Get().(*[32]byte)

			defer func() {
				intPool.Put(intHash)
				hashPool.Put(hash)
			}()

			for nonce := start; nonce < end; nonce++ {
				select {
				case <-stopChan:
					return
				default:
					data := pow.InitData(nonce)
					*hash = sha256.Sum256(data)
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

func (pow *ProofOfWork) Sha256Validate() bool {
	var intHash big.Int

	data := pow.InitData(pow.Block.Nonce)

	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:])

	return intHash.Cmp(pow.Target) == -1
}
