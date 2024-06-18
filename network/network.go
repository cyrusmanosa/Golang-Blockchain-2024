package network

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"golang-blockchain/blockchain"
	er "golang-blockchain/err"
	"log"
	"net"
	"os"
	"runtime"
	"syscall"

	"github.com/vrecan/death"
)

const (
	Protocol      = "tcp"
	Version       = 1
	CommandLength = 12
)

var (
	NodeAddress     string
	MinerAddress    string
	KnownNodes      = []string{"localhost:3000"}
	BlocksInTransit = [][]byte{}
	MemoryPool      = make(map[string]blockchain.Transaction)
)

type Addr struct {
	AddrList []string
}

type Block struct {
	AddrFrom string
	Block    []byte
}
type GetBlocks struct {
	AddrFrom string
}
type GetData struct {
	AddrFrom string
	Type     string
	ID       []byte
}
type Inv struct {
	AddrFrom string
	Type     string
	Items    [][]byte
}
type Tx struct {
	AddrForm    string
	Transaction []byte
}
type Versions struct {
	Version    int
	BestHeight int
	AddrFrom   string
}

func CmdToBytes(cmd string) []byte {
	var bytes [CommandLength]byte
	for i, c := range cmd {
		bytes[i] = byte(c)
	}
	return bytes[:]
}

func BytesToCmd(bytes []byte) string {
	var cmd []byte
	for _, b := range bytes {
		if b != 0x0 {
			cmd = append(cmd, b)
		}
	}
	return string(cmd)
}

func GobEncode(data interface{}) []byte {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(data)
	er.Handle(err)

	return buff.Bytes()
}

func CloseDB(chain *blockchain.BlockChain) {
	d := death.NewDeath(syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	d.WaitForDeathWithFunc(func() {
		defer os.Exit(1)
		defer runtime.Goexit()
		chain.Database.Close()
	})
}

func RequestBlocks() {
	for _, node := range KnownNodes {
		SendGetBlocks(node)
	}
}

func ExtractCmd(request []byte) []byte {
	return request[:CommandLength]
}

func MineTx(chain *blockchain.BlockChain) {
	var txs []*blockchain.Transaction

	for id := range MemoryPool {
		fmt.Printf("tx: %s\n", MemoryPool[id].ID)
		tx := MemoryPool[id]
		if chain.VerifyTransaction(&tx) {
			txs = append(txs, &tx)
		}
	}

	if len(txs) == 0 {
		fmt.Println("All Transaction are invalid")
		return
	}

	cbTx := blockchain.CoinbaseTx(MinerAddress, "")
	txs = append(txs, cbTx)

	newBlock := chain.MineBlock(txs)
	UTXOSet := blockchain.UTXOSet{Blockchain: chain}
	UTXOSet.Reindex()

	fmt.Println("New Block mined")

	for _, tx := range txs {
		txID := hex.EncodeToString(tx.ID)
		delete(MemoryPool, txID)
	}

	for _, node := range KnownNodes {
		if node != NodeAddress {
			SendInv(node, "block", [][]byte{newBlock.Hash})
		}
	}

	if len(MemoryPool) > 0 {
		MineTx(chain)
	}
}

func StartServer(nodeID, minerAddress string) {
	NodeAddress = fmt.Sprintf("localhost:%s", nodeID)
	MinerAddress = minerAddress
	ln, err := net.Listen(Protocol, NodeAddress)
	if err != nil {
		log.Panic(err)
	}
	defer ln.Close()

	chain := blockchain.ContinueBlockChain(nodeID)
	defer chain.Database.Close()
	go CloseDB(chain)

	if NodeAddress != KnownNodes[0] {
		SendVersions(KnownNodes[0], chain)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Panic(err)
		}
		go HandleConnection(conn, chain)

	}
}

func NodeIsKnown(addr string) bool {
	for _, node := range KnownNodes {
		if node == addr {
			return true
		}
	}
	return false
}
