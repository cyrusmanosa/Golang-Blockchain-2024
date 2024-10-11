package main

import (
	"blockchain-back/controllers"
	"blockchain-back/dsl"
	"blockchain-back/server"
	_ "net/http/pprof"
	"os"
	"runtime"
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
		server.Gin(data)
	default:
		runtime.Goexit()
	}
}
