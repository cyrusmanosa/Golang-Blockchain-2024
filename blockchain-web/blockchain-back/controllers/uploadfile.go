package controllers

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func UploadFile(c *gin.Context) {
	// 从表单中获取文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "获取文件失败"})
		return
	}

	// 生成保存文件的路径
	filename := filepath.Base(file.Filename)
	dst := filepath.Join("/Users/cyrusman/Desktop/ProgrammingLearning/GolangBlockchain2024/blockchain-web/blockchain-back/dsl/Original/", filename)

	// 保存文件
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("文件 '%s' 上传成功", filename),
	})
}
