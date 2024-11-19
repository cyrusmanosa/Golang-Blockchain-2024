package blockchain

import (
	"fmt"
	"log"
	"math"
	"math/big"

	"golang.org/x/crypto/argon2"
	"golang.org/x/exp/rand"
)

func Argon2CreateHash(data, salt []byte) ([]byte, error) {
	hash := argon2.IDKey(data, salt, 1, 8*1024, 6, 32)
	return hash, nil
}

func Argon2Salt() ([]byte, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, fmt.Errorf("Argon2 Salt have a error: %w", err)
	}
	return salt, nil
}

func (pow *ProofOfWork) Argon2Run() (int, []byte) {
	var intHash big.Int
	var argon2 []byte

	nonce := 0
	fmt.Println("Loading................")
	salt, err := Argon2Salt()
	if err != nil {
		log.Panic(err)
	}

	for nonce < math.MaxInt64 {
		data := pow.InitData(nonce)

		argon2, err = Argon2CreateHash(data, salt)
		if err != nil {
			log.Panic(err)
		}

		// fmt.Printf("\r%x", argon2)
		intHash.SetBytes(argon2[:])

		if intHash.Cmp(pow.Target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Println()

	return nonce, argon2[:]
}

func (pow *ProofOfWork) Argon2Validate() bool {
	var intHash big.Int

	data := pow.InitData(pow.Block.Nonce)
	salt, err := Argon2Salt()
	if err != nil {
		log.Panic(err)
	}

	argon2, err := Argon2CreateHash(data, salt)
	if err != nil {
		log.Panic(err)
	}
	intHash.SetBytes(argon2[:])

	return intHash.Cmp(pow.Target) == -1
}
