package blockchain

import (
	"fmt"

	"golang.org/x/crypto/argon2"
	"golang.org/x/exp/rand"
)

func Argon2CreateHash(data, salt []byte) ([]byte, error) {
	hash := argon2.IDKey(data, salt, 1, 4*1024, 1, 16)
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
