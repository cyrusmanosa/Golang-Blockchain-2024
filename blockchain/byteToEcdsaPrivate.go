package blockchain

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/x509"
	"encoding/pem"
	"math/big"
)

func BytesToECDSAPrivateKey(keyBytes []byte) (*ecdsa.PrivateKey, error) {
	// Attempt to decode PEM-encoded key
	block, _ := pem.Decode(keyBytes)
	if block != nil {
		keyBytes = block.Bytes
	}

	// Parse the key bytes
	key, err := x509.ParseECPrivateKey(keyBytes)
	if err == nil {
		return key, nil
	}

	// If not a PEM-encoded key, try to parse as a raw private key
	curve := elliptic.P256() // Change this if you are using a different curve
	priv := new(big.Int).SetBytes(keyBytes)
	return &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: curve,
			X:     nil,
			Y:     nil,
		},
		D: priv,
	}, nil
}
