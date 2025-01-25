package main

import (
	"os"
	"runtime"

	"blockchain-back/controllers"
	"blockchain-back/server"
)

func main() {
	if len(os.Args) < 2 {
		runtime.Goexit()
	}

	infPath := "/Users/cyrusman/Desktop/ProgrammingLearning/Golang-Blockchain-2024/blockchain-back/dsl/Original/"
	controllers.DeleteAllFilesInFolder(infPath)
	infPath2 := "/Users/cyrusman/Desktop/ProgrammingLearning/Golang-Blockchain-2024/blockchain-back/tmp"
	controllers.DeleteAllFilesInFolder(infPath2)
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
