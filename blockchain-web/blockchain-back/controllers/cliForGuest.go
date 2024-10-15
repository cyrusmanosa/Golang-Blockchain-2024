package controllers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	blockchain "blockchain-back/blockchain"
	"blockchain-back/dsl"
	models "blockchain-back/modules"

	"github.com/gin-gonic/gin"
)

type CommandLine struct {
	blockchain *blockchain.BlockChain
}

const (
	infPath = "/Users/cyrusman/Desktop/ProgrammingLearning/GolangBlockchain2024/blockchain-web/blockchain-back/dsl/Original/履歴書.pdf"
	outPath = "/Users/cyrusman/Desktop/ProgrammingLearning/GolangBlockchain2024/blockchain-web/blockchain-back/dsl/Svg/aaa.svg"
	layout  = "2006-01-02 15:04:05"
)

// / ----------------------------- Not Confirm or Text -----------------------------------
func AddBlockForGin(ctx *gin.Context) {
	chain := blockchain.InitBlockChainForGuest()
	defer chain.Database.Close()
	cli := CommandLine{chain}

	// take Data
	var req models.InputData
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

	chain := blockchain.InitBlockChainForGuest()
	defer chain.Database.Close()
	cli := CommandLine{chain}

	iter := cli.blockchain.Iterator()

	for {
		block := iter.Next()
		var prettyJSON bytes.Buffer
		err := json.Indent(&prettyJSON, block.Data, "", "    ")
		if err != nil {
			ErrorResponse(err)
			break
		}

		var dataMap map[string]interface{}
		err = json.Unmarshal(block.Data, &dataMap)
		if err != nil {
			ErrorResponse(err)
			break
		}

		dataName, ok := dataMap["name"].(string)
		if !ok {
			ErrorResponse(err)
			break
		}

		if dataName == rspName {
			dataEmail, _ := dataMap["email"].(string)
			dataCN, _ := dataMap["company_name"].(string)
			dataT, _ := dataMap["send_time"].(string)
			dataMsg, _ := dataMap["message"].(string)

			svgData, err := dsl.PdfToSvg(infPath, outPath)
			if err != nil {
				ErrorResponse(err)
				break
			} else {
				svgBase64 := base64.StdEncoding.EncodeToString(svgData)
				newData := models.InputData{
					Name:        dataName,
					Email:       dataEmail,
					CompanyName: dataCN,
					Message:     dataMsg,
					Cv:          svgBase64,
					Status:      "Checked",
					SendTime:    dataT,
					ConfirmTime: time.Now().Format(layout),
				}

				cli.blockchain.AddBlockForGuest(newData)
				fmt.Println("Added Block!")
				fmt.Println(" ")
				cli.PrintChainForConfirm()
				SendRsp(newData)
			}
			if len(block.PrevHash) == 0 {
				break
			}
		}
	}
}

// / ----------------------------- Error func -----------------------------------
func ErrorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
