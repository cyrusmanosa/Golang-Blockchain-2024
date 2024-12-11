package blockchain

import (
	"fmt"
	"log"
	"math"
	"math/big"
	"sync"

	"github.com/pedroalbanese/skein"
)

func (pow *ProofOfWork) SkeinLowRun() (int, []byte) {
	var intHash big.Int
	var hash []byte
	nonce := 0

	fmt.Println("\n-Low- Loading................")

	for nonce < math.MaxInt64 {
		data := pow.InitData(nonce)
		hasher := skein.New256(nil)

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
func (pow *ProofOfWork) SkeinRun() (int, []byte) {

	numCPUs := 4

	fmt.Println("\n-High- Loading................")

	rangeSize := math.MaxInt64 / numCPUs

	var resultNonce int
	var resultHash []byte
	stopChan := make(chan struct{})
	resultChan := make(chan struct {
		nonce int
		hash  []byte
	}, numCPUs)

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

	var wg sync.WaitGroup
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
					hasher := skein.New256(nil)
					hasher.Write(data)
					hash = hasher.Sum(hash[:0])
					intHash.SetBytes(hash)
					if intHash.Cmp(pow.Target) == -1 {
						select {
						case resultChan <- struct {
							nonce int
							hash  []byte
						}{nonce: nonce, hash: append([]byte(nil), hash...)}:
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

	result := <-resultChan
	close(stopChan)
	resultNonce = result.nonce
	resultHash = result.hash

	return resultNonce, resultHash
}
func (pow *ProofOfWork) SkeinValidate() bool {
	var intHash big.Int
	hasher := skein.New256(nil)

	data := pow.InitData(pow.Block.Nonce)

	if _, err := hasher.Write(data); err != nil {
		log.Panic("Error while hashing data: ", err)
	}

	hash := hasher.Sum(nil)

	intHash.SetBytes(hash)

	return intHash.Cmp(pow.Target) == -1
}
