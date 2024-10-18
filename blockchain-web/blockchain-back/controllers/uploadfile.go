package controllers

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get file"})
		return
	}

	filename := filepath.Base(file.Filename)
	dst := filepath.Join("/Users/cyrusman/Desktop/ProgrammingLearning/GolangBlockchain2024/blockchain-web/blockchain-back/dsl/Original/", filename)

	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("File '%s' uploaded successfully", filename),
	})
}
