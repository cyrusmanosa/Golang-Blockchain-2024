package main

import (
	"Argon2/controllers"
	"Argon2/server"
	_ "net/http/pprof"
	"os"
	"runtime"
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
