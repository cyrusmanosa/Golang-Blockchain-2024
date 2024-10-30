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
		dataCv, _ := dataMap["file"].(string)
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
			fmt.Println("	File: OK")
		} else {
			fmt.Println("	File: ")
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
	chain, err := blockchain.InitBlockChainForGuest()
	if err != nil {
		ErrorResponse(err)
		return nil
	}

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
		fileValue, fileExists := dataMap["file"]
		if !fileExists {
			log.Println("Field 'file' not found in dataMap")
			continue
		}

		if dataName == CusName {
			switch v := fileValue.(type) {
			case string:
				svgBase64, err := base64.StdEncoding.DecodeString(v)
				if err != nil {
					log.Println("Failed to decode base64:", err)
					return nil
				}
				return svgBase64
			case []byte:
				return v
			default:
				log.Println("Field 'file' is neither a string nor []byte")
				continue
			}
		} else {
			log.Println("!!!!!!! cv is empty !!!!!!!")
		}
	}
	return nil
}
