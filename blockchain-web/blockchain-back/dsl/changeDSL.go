package dsl

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func PdfToSvg(infPath string, outPath string) ([]byte, error) {
	files, err := os.ReadDir(infPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read directory: %s, error: %v", infPath, err)
	}

	file := files[0]
	filePath := filepath.Join(infPath, file.Name())

	if !strings.HasSuffix(strings.ToLower(file.Name()), ".pdf") {
		return nil, fmt.Errorf("file is not in PDF format: %s", file.Name())
	}

	var stderr bytes.Buffer
	cmd := exec.Command("pdf2svg", filePath, outPath, "1")

	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		log.Println("Error running pdf2svg: ", err)
		log.Println("pdf2svg stderr: ", stderr.String())
		return nil, err
	}

	svgBytes, err := os.ReadFile(outPath)
	if err != nil {
		log.Println("Error reading SVG file:", err)
		return nil, err
	}

	return svgBytes, nil
}
