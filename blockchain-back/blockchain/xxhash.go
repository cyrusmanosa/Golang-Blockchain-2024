package blockchain

import (
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
		hash := xxHash256(data)
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
					hash := xxHash256(data)
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

func uint64ToBytesAtXxhash(num uint64) []byte {
	return []byte{
		byte(num >> 56), byte(num >> 48), byte(num >> 40), byte(num >> 32),
		byte(num >> 24), byte(num >> 16), byte(num >> 8), byte(num),
	}
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
	hash128 := make([]byte, 16)

	copy(hash128[:8], uint64ToBytesAtXxhash(hash1))
	copy(hash128[8:], uint64ToBytesAtXxhash(hash2))

	var hash256 [32]byte
	copy(hash256[:16], hash128)
	copy(hash256[16:], hash128)

	return hash256
}
