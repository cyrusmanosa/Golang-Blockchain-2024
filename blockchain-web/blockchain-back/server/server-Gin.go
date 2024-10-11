package server

import (
	"blockchain-back/controllers"
	"blockchain-back/dsl"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Gin(in, out string, filePath []byte) {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	r := gin.Default()
	r.Use(cors.Default())

	svgData := dsl.PdfToSvg(in, out)

	// /----------------------------------------------------------------------------------------------------------------------------------------------------------------
	r.GET("/png", func(c *gin.Context) {
		html := ` 
		<!DOCTYPE html>
		<html>
		<head>
		<title>Image Display</title>
		</head>
		<body>
		<h1>Image</h1>
		<img src="/imageTest" />
		</body>
		</html>`
		// Write the HTML response
		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, html)
	})
	r.GET("/imageTest", func(c *gin.Context) {
		c.Header("Content-Type", "image/png")
		c.Data(http.StatusOK, "image/png", filePath)
	})
	// /----------------------------------------------------------------------------------------------------------------------------------------------------------------

	r.GET("/svg/:Cname", func(c *gin.Context) {
		Cname := c.Param("Cname")
		Svg := controllers.TakeBlock(c, Cname)
		c.Header("Content-Type", "svg+xml")
		c.Data(http.StatusOK, "svg+xml", Svg)
	})

	r.GET("/image", func(c *gin.Context) {
		c.Header("Content-Type", "svg+xml")
		c.Data(http.StatusOK, "svg+xml", svgData)
	})
	r.PUT("/take", controllers.AddBlockForGin)
	r.PUT("/Check/:name", controllers.AddBlockForGinConfirm)

	fmt.Println("Starting server at :8080")
	log.Fatal(r.Run(":8080"))
}
