package blockchain

import (
	"bytes"
	"encoding/gob"
	er "golang-blockchain/err"
	"golang-blockchain/wallet"
)

type TxOutput struct {
	Value      int
	PubKeyHash []byte
}
type TxOutputs struct {
	Outputs []TxOutput
}

type TxInput struct {
	ID        []byte
	Out       int
	Signature []byte
	PubKey    []byte
}

func (outs TxOutputs) Serialize() []byte {
	var buffer bytes.Buffer
	encode := gob.NewEncoder(&buffer)
	err := encode.Encode(outs)
	er.Handle(err)
	return buffer.Bytes()
}

func DeserializeOutputs(data []byte) TxOutputs {
	var outputs TxOutputs
	decode := gob.NewDecoder(bytes.NewReader(data))
	err := decode.Decode(&outputs)
	er.Handle(err)
	return outputs
}

func (in *TxInput) UsesKey(PubKeyHash []byte) bool {
	lockingHash := wallet.PublicKeyHash(in.PubKey)
	return bytes.Equal(lockingHash, PubKeyHash)
}

func (out *TxOutput) Lock(address []byte) {
	PubKeyHash := wallet.Base58Decode(address)
	PubKeyHash = PubKeyHash[1 : len(PubKeyHash)-4]
	out.PubKeyHash = PubKeyHash
}

func (out *TxOutput) IsLockedWithKey(PubKeyHash []byte) bool {
	return bytes.Equal(out.PubKeyHash, PubKeyHash)
}

func NewTXOutput(value int, address string) *TxOutput {
	txo := &TxOutput{value, nil}
	txo.Lock([]byte(address))

	return txo
}
