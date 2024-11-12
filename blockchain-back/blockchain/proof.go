package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/big"
)

const Difficulty = 18

type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

func NewProof(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-Difficulty))

	pow := &ProofOfWork{b, target}

	return pow
}

func (pow *ProofOfWork) InitData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.Block.PrevHash,
			pow.Block.Data,
			ToHex(int64(nonce)),
			ToHex(int64(Difficulty)),
		},
		[]byte{},
	)

	return data
}

// /----------------------------- *********** -----------------------------------
func (pow *ProofOfWork) Sha256Run() (int, []byte) {
	var intHash big.Int
	var hash [32]byte

	nonce := 0
	fmt.Println("Loading................")
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
	fmt.Println()

	return nonce, hash[:]
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

///----------------------------- *********** -----------------------------------

func (pow *ProofOfWork) Sha256Validate() bool {
	var intHash big.Int

	data := pow.InitData(pow.Block.Nonce)

	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:])

	return intHash.Cmp(pow.Target) == -1
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

func ToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)

	}

	return buff.Bytes()
}
