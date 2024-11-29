package controllers

import (
	"fmt"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"

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

	files, err := t.UploadOneFile(c.Request, "/Users/cyrusman/Desktop/ProgrammingLearning/Golang-Blockchain-2024/blockchain-back/dsl/Original")
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	c.String(http.StatusOK, fmt.Sprintf("Uploaded 1 file, %s, to the uploads folder", files.OriginalFileName))
}

func GetUniquePDF(folderPath string) (string, error) {
	var pdfFiles []string

	err := filepath.WalkDir(folderPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && strings.HasSuffix(strings.ToLower(d.Name()), ".pdf") {
			pdfFiles = append(pdfFiles, path)
		}

		return nil
	})
	if err != nil {
		return "", err
	}

	if len(pdfFiles) == 0 {
		return "", fmt.Errorf("資料夾內沒有找到 PDF 文件")
	}
	if len(pdfFiles) > 1 {
		return "", fmt.Errorf("資料夾內找到多個 PDF 文件，請確保只有一個")
	}

	return pdfFiles[0], nil
}
