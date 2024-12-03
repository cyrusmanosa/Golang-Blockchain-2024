package blockchain

import (
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
	"runtime"
	"sync"

	"github.com/dgryski/go-farm"
)

func (pow *ProofOfWork) FarmLowRun() (int, []byte) {
	var intHash big.Int
	var hash [32]byte

	nonce := 0
	runtime.GC()
	fmt.Println("\n-Low- Loading................")

	for nonce < math.MaxInt64 {
		data := pow.InitData(nonce)
		hash := FarmHash(data)
		intHash.SetBytes(hash[:])

		if intHash.Cmp(pow.Target) == -1 {
			break
		}
		nonce++
	}

	return nonce, hash[:]
}

func (pow *ProofOfWork) FarmRun() (int, []byte) {
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
					hash := FarmHash(data)
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

func (pow *ProofOfWork) FarmValidate() bool {
	var intHash big.Int

	data := pow.InitData(pow.Block.Nonce)
	hash := FarmHash(data)

	intHash.SetBytes(hash[:])

	return intHash.Cmp(pow.Target) == -1
}

func FarmHash(data []byte) [32]byte {
	hash64_1 := farm.Hash64(data)
	hash64_2 := farm.Hash64(append(data, '1'))
	hash64_3 := farm.Hash64(append(data, '2'))
	hash64_4 := farm.Hash32(append(data, '3'))

	var hash32Bytes [32]byte
	binary.LittleEndian.PutUint64(hash32Bytes[:8], hash64_1)
	binary.LittleEndian.PutUint64(hash32Bytes[8:16], hash64_2)
	binary.LittleEndian.PutUint64(hash32Bytes[16:24], hash64_3)
	binary.LittleEndian.PutUint32(hash32Bytes[24:], hash64_4)

	return hash32Bytes
}
