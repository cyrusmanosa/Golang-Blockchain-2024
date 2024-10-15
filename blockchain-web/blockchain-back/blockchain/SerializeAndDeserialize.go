package blockchain

import (
	"bytes"
	"encoding/gob"
	"log"
)

func (b *Block) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)
	err := encoder.Encode(b)
	if err != nil {
		log.Println("Serialize Error: ", err)
	}
	return res.Bytes()
}

func Deserialize(data []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	if err != nil {
		log.Println("Deserialize Error: ", err)
	}
	return &block
}
