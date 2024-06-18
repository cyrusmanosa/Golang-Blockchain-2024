package network

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"golang-blockchain/blockchain"
	er "golang-blockchain/err"
	"io"
	"net"
)

func HandleAddr(request []byte) {
	var buff bytes.Buffer
	var payload Addr

	buff.Write(request[CommandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	er.Handle(err)

	KnownNodes = append(KnownNodes, payload.AddrList...)
	fmt.Printf("there are %d known nodes\n", len(KnownNodes))
	RequestBlocks()
}

func HandleBlock(request []byte, chain *blockchain.BlockChain) {
	var buff bytes.Buffer
	var payload Block
	buff.Write(request[CommandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	er.Handle(err)

	blockData := payload.Block
	block := blockchain.Deserialize(blockData)
	fmt.Println("Recevied a new block")

	chain.AddBlock(block)
	fmt.Printf("Added block %x\n", block.Hash)

	if len(BlocksInTransit) > 0 {
		blockHash := BlocksInTransit[0]
		SendGetData(payload.AddrFrom, "block", blockHash)
		BlocksInTransit = BlocksInTransit[1:]
	} else {
		UTXOSet := blockchain.UTXOSet{Blockchain: chain}
		UTXOSet.Reindex()
	}
}

func HandleGetBlocks(request []byte, chain *blockchain.BlockChain) {
	var buff bytes.Buffer
	var payload GetBlocks

	buff.Write(request[CommandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	er.Handle(err)

	blocks := chain.GetBlockHashes()
	SendInv(payload.AddrFrom, "block", blocks)
}

func HandleGetData(request []byte, chain *blockchain.BlockChain) {
	var buff bytes.Buffer
	var payload GetData

	buff.Write(request[CommandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	er.Handle(err)

	if payload.Type == "block" {
		block, err := chain.GetBlock([]byte(payload.ID))
		er.Handle(err)
		SendBlock(payload.AddrFrom, &block)
	}
	if payload.Type == "tx" {
		txID := hex.EncodeToString(payload.ID)
		tx := MemoryPool[txID]

		SendTx(payload.AddrFrom, &tx)
	}
}

func HandleTx(request []byte, chain *blockchain.BlockChain) {
	var buff bytes.Buffer
	var payload Tx

	buff.Write(request[CommandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	er.Handle(err)

	txData := payload.Transaction
	tx := blockchain.DeserializeTransaction(txData)
	MemoryPool[hex.EncodeToString(tx.ID)] = tx

	fmt.Printf("%s,%d\n", NodeAddress, len(MemoryPool))
	if NodeAddress == KnownNodes[0] {
		for _, node := range KnownNodes {
			if node != NodeAddress && node != payload.AddrForm {
				SendInv(node, "tx", [][]byte{tx.ID})
			}
		}
	} else {
		if len(MemoryPool) >= 2 && len(MinerAddress) > 0 {
			MineTx(chain)
		}
	}
}

func HandleInv(request []byte, chain *blockchain.BlockChain) {
	var buff bytes.Buffer
	var payload Inv

	buff.Write(request[CommandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	er.Handle(err)

	fmt.Printf("Recevied inventory with %d %s\n", len(payload.Items), payload.Type)

	if payload.Type == "block" {
		BlocksInTransit = payload.Items
		blockHash := payload.Items[0]

		SendGetData(payload.AddrFrom, "block", blockHash)

		newInTransit := [][]byte{}
		for _, b := range BlocksInTransit {
			if !bytes.Equal(b, blockHash) {
				newInTransit = append(newInTransit, b)
			}
		}
		BlocksInTransit = newInTransit
	}

	if payload.Type == "tx" {
		txID := payload.Items[0]
		if MemoryPool[hex.EncodeToString(txID)].ID == nil {
			SendGetData(payload.AddrFrom, "tx", txID)
		}
	}
}

func HandleVersion(request []byte, chain *blockchain.BlockChain) {
	var buff bytes.Buffer
	var payload Versions

	buff.Write(request[CommandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	er.Handle(err)
	bestHeight := chain.GetBestHeight()
	otherHeight := payload.BestHeight
	if bestHeight < otherHeight {
		SendGetBlocks(payload.AddrFrom)
	} else if bestHeight > otherHeight {
		SendVersions(payload.AddrFrom, chain)
	}

	if !NodeIsKnown(payload.AddrFrom) {
		KnownNodes = append(KnownNodes, payload.AddrFrom)
	}
}

func HandleConnection(conn net.Conn, chain *blockchain.BlockChain) {
	req, err := io.ReadAll(conn)
	defer conn.Close()
	er.Handle(err)

	command := BytesToCmd(req[:CommandLength])
	fmt.Printf("Received %s command \n", command)

	switch command {
	case "addr":
		HandleAddr(req)
	case "block":
		HandleBlock(req, chain)
	case "inv":
		HandleInv(req, chain)
	case "getblocks":
		HandleGetBlocks(req, chain)
	case "getdata":
		HandleGetData(req, chain)
	case "tx":
		HandleTx(req, chain)
	case "version":
		HandleVersion(req, chain)
	default:
		fmt.Println("Unknown command")
	}
}
