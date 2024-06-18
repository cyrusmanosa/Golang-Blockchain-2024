package blockchain

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"math/big"

	"github.com/btcsuite/btcd/btcec/v2"
)

func BytesToECDSAPrivateKey(keyBytes []byte) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode(keyBytes)
	if block != nil {
		keyBytes = block.Bytes
	}

	key, err := x509.ParseECPrivateKey(keyBytes)
	if err == nil {
		return key, nil
	}

	curve := btcec.S256()
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
