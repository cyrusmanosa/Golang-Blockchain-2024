package blockchain

import (
	"context"
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

	fmt.Println("\n-Low- Loading................")

	for nonce < math.MaxInt64 {
		data := pow.InitData(nonce)
		hash = FarmHash(data)
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
					hash := FarmHash(data)
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

	result := <-resultChan
	return result.nonce, result.hash
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
