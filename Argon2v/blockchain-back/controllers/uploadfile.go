package controllers

import (
	"fmt"
	"net/http"

	toolkit "github.com/cyrusmanosa/Toolkit/v2"
	"github.com/gin-gonic/gin"
)

func UploadOneFiles(c *gin.Context) {
	t := toolkit.Tools{
		MaxFileSize: 1024 * 1024 * 1024,
		AllowedFileTypes: []string{
			"application/pdf",
			"image/jpeg",
			"image/png",
		},
	}

	files, err := t.UploadOneFile(c.Request, "/Users/cyrusman/Desktop/ProgrammingLearning/Golang-Blockchain-2024/Argon2v/blockchain-back/dsl/Original")
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	c.String(http.StatusOK, fmt.Sprintf("Uploaded 1 file, %s, to the uploads folder", files.OriginalFileName))
}
