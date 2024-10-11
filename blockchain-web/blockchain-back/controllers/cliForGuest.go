package controllers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
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
	pdfPath = "/Users/cyrusman/Desktop/ECCコンピューター専門学校/Year-3/Y3-Sem1/ITゼミ演習１/blockchain-web/blockchain-back/dsl/Original/履歴書.pdf"
	svgPath = "/Users/cyrusman/Desktop/ECCコンピューター専門学校/Year-3/Y3-Sem1/ITゼミ演習１/blockchain-web/blockchain-back/dsl/Original/Svg/TestSVG.svg"
	layout  = "2006-01-02 15:04:05"
)

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

func AddBlockForGinConfirm(ctx *gin.Context) {
	rspName := ctx.Param("name")

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
		var prettyJSON bytes.Buffer
		err := json.Indent(&prettyJSON, block.Data, "", "    ")
		if err != nil {
			ErrorResponse(err)
		}

		var dataMap map[string]interface{}
		err = json.Unmarshal(block.Data, &dataMap)
		if err != nil {
			ErrorResponse(err)
			return
		}
		dataName, ok := dataMap["name"].(string)
		if !ok {
			ErrorResponse(err)
			return
		}
		if dataName == rspName {
			dataEmail, _ := dataMap["email"].(string)
			dataCN, _ := dataMap["company_name"].(string)
			dataT, _ := dataMap["send_time"].(string)
			dataMsg, _ := dataMap["message"].(string)

			svgData := dsl.PdfToSvg(pdfPath, svgPath)
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

func ErrorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
