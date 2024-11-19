package main

import (
	_ "net/http/pprof"
	"os"
	"runtime"

	"blockchain-back/controllers"
	"blockchain-back/server"
)

func main() {
	if len(os.Args) < 2 {
		runtime.Goexit()
	}
	switch os.Args[1] {
	case "print":
		controllers.Cli()
	case "server":
		server.Gin()
	default:
		runtime.Goexit()
	}
}

// 2024-11-08
