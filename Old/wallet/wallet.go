package wallet

import (
	"bytes"
	"crypto/sha256"
	"log"

	"github.com/btcsuite/btcd/btcec/v2"
	"golang.org/x/crypto/ripemd160"
)

const (
	checksumLength = 4
	version        = byte(0x00)
	walletFile     = "./tmp/wallets.data"
)

type Wallet struct {
	PrivateKey []byte
	PublicKey  []byte
}

func (w Wallet) Address() []byte {
	pubHash := PublicKeyHash(w.PublicKey)

	versionedHash := append([]byte{version}, pubHash...)
	checksum := Checksum(versionedHash)

	fullHash := append(versionedHash, checksum...)
	address := Base58Encode(fullHash)

	return address
}

func NewKeyPair() ([]byte, []byte) {
	private, err := btcec.NewPrivateKey()
	if err != nil {
		log.Panic(err)
	}

	privKeyBytes := private.Serialize()
	pubKeyBytes := private.PubKey().SerializeCompressed()

	return privKeyBytes, pubKeyBytes
}

func MakeWallet() *Wallet {
	private, public := NewKeyPair()
	wallet := Wallet{PrivateKey: private, PublicKey: public}

	return &wallet
}

func PublicKeyHash(pubKey []byte) []byte {
	pubHash := sha256.Sum256(pubKey)

	hasher := ripemd160.New()
	_, err := hasher.Write(pubHash[:])
	if err != nil {
		log.Panic(err)
	}

	publicRipMD := hasher.Sum(nil)

	return publicRipMD
}

func Checksum(payload []byte) []byte {
	firstHash := sha256.Sum256(payload)
	secondHash := sha256.Sum256(firstHash[:])

	return secondHash[:checksumLength]
}

func ValidateAddress(address string) bool {
	PubKeyHash := Base58Decode([]byte(address))
	actualChecksum := PubKeyHash[len(PubKeyHash)-checksumLength:]
	version := PubKeyHash[0]
	PubKeyHash = PubKeyHash[1 : len(PubKeyHash)-checksumLength]
	targetChecksum := Checksum(append([]byte{version}, PubKeyHash...))

	return bytes.Compare(actualChecksum, targetChecksum) == 0
}
