package controllers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"blockchain-back/blockchain"
	"blockchain-back/dsl"
	"blockchain-back/modules"
)

type CommandLine struct {
	blockchain *blockchain.BlockChain
}

const (
	infPath = "/Users/cyrusman/Desktop/ProgrammingLearning/Golang-Blockchain-2024/blockchain-back/dsl/Original/"
	outPath = "/Users/cyrusman/Desktop/ProgrammingLearning/Golang-Blockchain-2024/blockchain-back/dsl/Svg/"
	layout  = "2006-01-02 15:04:05"
)

// / ----------------------------- Not Confirm or Text -----------------------------------
func AddBlockForGin(ctx *gin.Context) {
	chain, err := blockchain.InitBlockChainForGuest()
	if err != nil {
		ErrorResponse(err)
	}
	defer chain.Database.Close()

	cli := CommandLine{chain}

	// take Data
	var req modules.InputData
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
		log.Println(req)
		return
	}

	req.SendTime = time.Now().Format(layout)
	req.Status = "Unconfirmed"

	cli.blockchain.AddBlockForGuest(req)
	fmt.Println("Added Block!")
	SendRequest(req)
	fmt.Println(" ")
	cli.PrintChain()
}

// / ----------------------------- Confirm -----------------------------------
func AddBlockForGinConfirm(ctx *gin.Context) {
	rspName := ctx.Param("name")

	chain, err := blockchain.InitBlockChainForGuest()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}
	defer chain.Database.Close()

	cli := CommandLine{chain}
	iter := cli.blockchain.Iterator()

	for {
		block := iter.Next()

		if block == nil {
			ctx.JSON(http.StatusInternalServerError, ErrorResponse(fmt.Errorf("block is nil")))
			return
		}

		var prettyJSON bytes.Buffer
		err := json.Indent(&prettyJSON, block.Data, "", "    ")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
			return
		}

		var dataMap map[string]interface{}
		err = json.Unmarshal(block.Data, &dataMap)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
			return
		}

		dataName, ok := dataMap["name"].(string)
		if !ok {
			ctx.JSON(http.StatusInternalServerError, ErrorResponse(fmt.Errorf("name field not found")))
			return
		}

		if dataName == rspName {
			dataEmail, _ := dataMap["email"].(string)
			dataCN, _ := dataMap["company_name"].(string)
			dataT, _ := dataMap["send_time"].(string)
			dataMsg, _ := dataMap["message"].(string)
			dataHash, _ := dataMap["hash"].(string)

			r := RandomString()
			outPath2 := fmt.Sprint(outPath, r, ".svg")
			svgData, err := dsl.PdfToSvg(infPath, outPath2)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
				return
			}
			svgBase64 := base64.StdEncoding.EncodeToString(svgData)

			newData := modules.InputData{
				Name:        dataName,
				Email:       dataEmail,
				CompanyName: dataCN,
				Message:     dataMsg,
				Hash:        dataHash,
				File:        svgBase64,
				Status:      "Checked",
				SendTime:    dataT,
			}

			cli.blockchain.AddBlockForGuest(newData)
			fmt.Println("Added Block!")

			var wg sync.WaitGroup
			wg.Add(2)
			go func() {
				defer wg.Done()
				cli.PrintChain()
			}()
			go func() {
				defer wg.Done()
				SendRsp(newData)
			}()
			wg.Wait()

			// delete pdf and svg
			err = DeleteAllFilesInFolder(infPath)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
				return
			}
			err = DeleteAllFilesInFolder(outPath)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
				return
			}
			return
		}
		if len(block.PrevHash) == 0 {
			break
		}
	}
}

// / ----------------------------- Error func -----------------------------------
func ErrorResponse(err error) gin.H {
	log.Println("Error: ", err)
	return gin.H{"error": err.Error()}
}
