package controllers

import (
	"blockchain-back/blockchain"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (cli *CommandLine) PrintChain() {
	iter := cli.blockchain.Iterator()
	for {
		block := iter.Next()
		if len(block.PrevHash) == 0 {
			break
		}
		fmt.Println("************************************************************************************************")
		if block.PrevHash != nil {
			fmt.Printf("Prev. hash: %x\n", block.PrevHash)
		}
		var prettyJSON bytes.Buffer
		err := json.Indent(&prettyJSON, block.Data, "", "    ")
		if err != nil {
			log.Println("Failed to format JSON:", err)
			break
		}

		var dataMap map[string]interface{}
		err = json.Unmarshal(block.Data, &dataMap)
		if err != nil {
			log.Print("Failed to parse JSON:", err)
			break
		}

		dataName, _ := dataMap["name"].(string)
		dataEmail, _ := dataMap["email"].(string)
		dataCN, _ := dataMap["company_name"].(string)
		dataCv, _ := dataMap["cv"].(string)
		dataMsg, _ := dataMap["message"].(string)
		dataStatus, _ := dataMap["status"].(string)
		dataST, _ := dataMap["send_time"].(string)
		dataCT, _ := dataMap["confirm_time"].(string)

		fmt.Println("Data:")
		fmt.Println("	Name: ", dataName)
		fmt.Println("	Email: ", dataEmail)
		fmt.Println("	CompanyName: ", dataCN)
		fmt.Println("	Message: ", dataMsg)
		if dataCv != "" || len(dataCv) > 100 {
			fmt.Println("	Cv: OK")
		} else {
			fmt.Println("	Cv: ")
		}
		fmt.Println("	Status: ", dataStatus)
		fmt.Println("	SendTime: ", dataST)
		fmt.Println("	ConfirmTime: ", dataCT)

		fmt.Printf("Hash: %x\n", block.Hash)
		pow := blockchain.NewProof(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
	}
}

func (cli *CommandLine) PrintChainForConfirm() {
	iter := cli.blockchain.Iterator()
	for {
		block := iter.Next()
		if len(block.PrevHash) == 0 {
			break
		}

		fmt.Println("************************************************************************************************")
		if block.PrevHash != nil {
			fmt.Printf("Prev. hash: %x\n", block.PrevHash)
		}
		var prettyJSON bytes.Buffer
		err := json.Indent(&prettyJSON, block.Data, "", "    ")
		if err != nil {
			log.Println("Failed to format JSON:", err)
			break
		}

		var dataMap map[string]interface{}
		err = json.Unmarshal(block.Data, &dataMap)
		if err != nil {
			log.Print("Failed to parse JSON:", err)
			break
		}
		dataName, _ := dataMap["name"].(string)
		dataEmail, _ := dataMap["email"].(string)
		dataCN, _ := dataMap["company_name"].(string)
		dataCv, _ := dataMap["cv"].(string)
		dataMsg, _ := dataMap["message"].(string)
		dataStatus, _ := dataMap["status"].(string)
		dataST, _ := dataMap["send_time"].(string)
		dataCT, _ := dataMap["confirm_time"].(string)

		fmt.Println("Data:")
		fmt.Println("	Name: ", dataName)
		fmt.Println("	Email: ", dataEmail)
		fmt.Println("	CompanyName: ", dataCN)
		fmt.Println("	Message: ", dataMsg)
		if dataCv != "" {
			fmt.Println("	Cv: OK")
		} else {
			fmt.Println("	Cv: ")
		}
		fmt.Println("	Status: ", dataStatus)
		fmt.Println("	SendTime: ", dataST)
		fmt.Println("	ConfirmTime: ", dataCT)

		fmt.Printf("Hash: %x\n", block.Hash)
		pow := blockchain.NewProof(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
	}
}

func TakeBlock(c *gin.Context, CusName string) []byte {
	chain := blockchain.InitBlockChainForGuest()
	defer chain.Database.Close()
	cli := CommandLine{chain}

	if _, err := os.Stat("./tmp/block"); os.IsNotExist(err) {
		err := os.MkdirAll("./tmp/block", os.ModePerm)
		if err != nil {
			log.Panic("Error Creating Dir: ", err)
		}
	}

	iter := cli.blockchain.Iterator()
	for {
		block := iter.Next()
		if len(block.PrevHash) == 0 {
			break
		}

		fmt.Println("************************************************************************************************")
		if block.PrevHash != nil {
			fmt.Printf("Prev. hash: %x\n", block.PrevHash)
		}

		var prettyJSON bytes.Buffer
		err := json.Indent(&prettyJSON, block.Data, "", "    ")
		if err != nil {
			fmt.Println("Failed to format JSON:", err)
			return nil
		}

		var dataMap map[string]interface{}
		err = json.Unmarshal(block.Data, &dataMap)
		if err != nil {
			log.Print("Failed to parse JSON:", err)
			return nil
		}

		dataName, _ := dataMap["name"].(string)
		dataCv, _ := dataMap["cv"].(string)
		if dataName == CusName {
			if dataCv != "" {
				svgBase64, _ := base64.StdEncoding.DecodeString(dataCv)
				return svgBase64
			} else {
				log.Println("!!!!!!! cv is empty !!!!!!!")
				os.Exit(1)
			}
		}
	}
	return nil
}
