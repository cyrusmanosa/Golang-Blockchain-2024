package dsl

import (
	"bytes"
	"log"
	"os"
	"os/exec"
)

func PdfToSvg(infPath string, outPath string) ([]byte, error) {
	var stderr bytes.Buffer
	cmd := exec.Command("pdf2svg", infPath, outPath, "1")
	cmd.Stderr = &stderr
	err := cmd.Run()
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
