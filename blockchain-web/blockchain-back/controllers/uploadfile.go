package controllers

import (
	"fmt"
	"net/http"

	toolkit "github.com/cyrusmanosa/Toolkit/v2"
	"github.com/gin-gonic/gin"
)

func UploadFiles(c *gin.Context) {

	t := toolkit.Tools{
		MaxFileSize:      1024 * 1024 * 1024, // 1 GB
		AllowedFileTypes: []string{"application/pdf"},
	}

	files, err := t.UploadFiles(c.Request, "/Users/cyrusman/Desktop/ProgrammingLearning/Udemy/Toolkit/appTest/UploadTest/uploads")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out := ""
	for _, item := range files {
		out += fmt.Sprintf("Uploaded %s to the uploads folder as %s\n", item.OriginalFileName, item.NewFileName)
	}

	c.String(http.StatusOK, out)
}
