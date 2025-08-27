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

	UploadFilePath := "/Users/cyrusman/Desktop/ProgrammingLearning/Project/Golang-Blockchain-2024/blockchain-back/dsl/Original/"
	controllers.DeleteAllFilesInFolder(UploadFilePath)

	ChainDataPath := "/Users/cyrusman/Desktop/ProgrammingLearning/Project/Golang-Blockchain-2024/blockchain-back/tmp"
	controllers.DeleteAllFilesInFolder(ChainDataPath)

	switch os.Args[1] {
	case "print":
		controllers.Cli()
	case "server":
		server.Gin()
	default:
		runtime.Goexit()
	}
}
