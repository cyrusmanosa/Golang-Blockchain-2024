package dsl

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"os"
	"os/exec"
)

func PngToSvg() []byte {
	path := "/Users/cyrusman/Desktop/ECCコンピューター専門学校/Year-3/Y3-Sem1/ITゼミ演習１/blockchain-web/blockchain-back/dsl/Original/TestPng.png"
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("Error decoding image:", err)
		return nil
	}

	var buf bytes.Buffer
	err = png.Encode(&buf, img)
	if err != nil {
		fmt.Println("Error encoding image to PNG:", err)
		return nil
	}
	pngData := buf.Bytes()

	return pngData
}

func PdfToSvg(pdfPath string, svgPath string) []byte {
	var stderr bytes.Buffer
	cmd := exec.Command("pdf2svg", pdfPath, svgPath, "1")
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error running pdf2svg: %v\n", err)
		fmt.Printf("pdf2svg stderr: %s\n", stderr.String())
		return nil
	}

	svgBytes, err := os.ReadFile(svgPath)
	if err != nil {
		fmt.Println("Error reading SVG file:", err)
		return nil
	}

	return svgBytes
}
