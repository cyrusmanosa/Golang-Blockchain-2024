package main

import (
	"blockchain-back/controllers"
	"blockchain-back/dsl"
	"blockchain-back/server"
	_ "net/http/pprof"
	"os"
	"runtime"
)

const (
	pdfPath = "/Users/cyrusman/Desktop/ECCコンピューター専門学校/Year-3/Y3-Sem1/ITゼミ演習１/blockchain-web/blockchain-back/dsl/Original/履歴書.pdf"
	svgPath = "/Users/cyrusman/Desktop/ECCコンピューター専門学校/Year-3/Y3-Sem1/ITゼミ演習１/blockchain-web/blockchain-back/dsl/Original/Svg/TestSVG.svg"
)

func main() {
	if len(os.Args) < 2 {
		runtime.Goexit()
	}

	switch os.Args[1] {
	case "command":
		controllers.Cli()
	case "server":
		data := dsl.PngToSvg()
		server.Gin(pdfPath, svgPath, data)
	default:
		runtime.Goexit()
	}
}
