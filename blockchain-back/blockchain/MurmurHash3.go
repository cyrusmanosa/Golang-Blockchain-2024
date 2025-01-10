package blockchain

import (
	"fmt"
	"math"
	"math/big"
	"runtime"
	"sync"

	"github.com/twmb/murmur3"
)

func (pow *ProofOfWork) MurmurHashLowRun() (int, []byte) {
	var intHash big.Int
	var hash [32]byte

	nonce := 0

	fmt.Println("\n-Low- Loading................")

	for nonce < math.MaxInt64 {
		data := pow.InitData(nonce)
		hash = MurmurHash256(data)
		intHash.SetBytes(hash[:])

		if intHash.Cmp(pow.Target) == -1 {
			break
		}
		nonce++
	}

	return nonce, hash[:]
}
func (pow *ProofOfWork) MurmurHashRun() (int, []byte) {
	numCPUs := 4
	runtime.GOMAXPROCS(numCPUs)
	fmt.Println("\n-High- Loading................")

	var resultNonce int
	var resultHash []byte
	stopChan := make(chan struct{})
	resultChan := make(chan struct {
		nonce int
		hash  []byte
	}, numCPUs)

	rangeSize := math.MaxInt64 / numCPUs
	var wg sync.WaitGroup

	intPool := &sync.Pool{
		New: func() interface{} {
			return new(big.Int)
		},
	}
	hashPool := &sync.Pool{
		New: func() interface{} {
			return make([]byte, 32)
		},
	}

	for i := 0; i < numCPUs; i++ {
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()

			intHash := intPool.Get().(*big.Int)
			defer intPool.Put(intHash)

			hash := hashPool.Get().([]byte)
			defer hashPool.Put(hash)

			for nonce := start; nonce < end; nonce++ {
				select {
				case <-stopChan:
					return
				default:
					data := pow.InitData(nonce)
					hash := MurmurHash256(data)
					intHash.SetBytes(hash[:])

					if intHash.Cmp(pow.Target) == -1 {
						select {
						case resultChan <- struct {
							nonce int
							hash  []byte
						}{nonce: nonce, hash: hash[:]}:
							close(stopChan)
							return
						case <-stopChan:
							return
						}
					}
				}
			}
		}(i*rangeSize, (i+1)*rangeSize)
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

func uint64ToBytesAtMurmur(n uint64) []byte {
	buf := make([]byte, 8)
	for i := uint(0); i < 8; i++ {
		buf[i] = byte(n >> (56 - i*8))
	}
	return buf
}

func (pow *ProofOfWork) MurmurHashValidate() bool {
	var intHash big.Int

	data := pow.InitData(pow.Block.Nonce)
	hash := MurmurHash256(data)

	intHash.SetBytes(hash[:])

	return intHash.Cmp(pow.Target) == -1
}

func MurmurHash256(data []byte) [32]byte {
	h1 := murmur3.SeedSum64(0, data)
	h2 := murmur3.SeedSum64(1, data)
	h3 := murmur3.SeedSum64(2, data)
	h4 := murmur3.SeedSum64(3, data)

	hash := [32]byte{}
	copy(hash[:8], uint64ToBytesAtMurmur(h1))
	copy(hash[8:16], uint64ToBytesAtMurmur(h2))
	copy(hash[16:24], uint64ToBytesAtMurmur(h3))
	copy(hash[24:], uint64ToBytesAtMurmur(h4))

	return hash
}
