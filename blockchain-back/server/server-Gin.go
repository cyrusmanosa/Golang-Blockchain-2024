package server

import (
	"log"
	"net/http"
	"os"
	"runtime/trace"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"blockchain-back/controllers"
)

func Gin() {
	traceFile, err := os.Create("trace.out")
	if err != nil {
		log.Fatalf("failed to create trace file: %v", err)
	}
	defer traceFile.Close()

	if err := trace.Start(traceFile); err != nil {
		log.Fatalf("failed to start trace: %v", err)
	}
	defer trace.Stop()

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	r := gin.Default()
	r.Use(cors.Default())

	r.GET("/svg/:Cname", func(c *gin.Context) {
		trace.Log(c.Request.Context(), "Handler", "/svg/:Cname - Start")
		Cname := c.Param("Cname")
		Svg := controllers.TakeBlock(c, Cname)
		c.Header("Content-Type", "svg+xml")
		c.Data(http.StatusOK, "svg+xml", Svg)
		trace.Log(c.Request.Context(), "Handler", "/svg/:Cname - End")
	})

	r.POST("/Upload", func(c *gin.Context) {
		trace.Log(c.Request.Context(), "Handler", "/Upload - Start")
		controllers.UploadOneFiles(c)
		trace.Log(c.Request.Context(), "Handler", "/Upload - End")
	})

	r.POST("/take", func(c *gin.Context) {
		trace.Log(c.Request.Context(), "Handler", "/take - Start")
		controllers.AddBlockForGin(c)
		trace.Log(c.Request.Context(), "Handler", "/take - End")
	})

	r.POST("/Check/:name", func(c *gin.Context) {
		trace.Log(c.Request.Context(), "Handler", "/Check/:name - Start")
		controllers.AddBlockForGinConfirm(c)
		trace.Log(c.Request.Context(), "Handler", "/Check/:name - End")
	})

	log.Println("Starting server at :8080")
	r.Run(":8080")
}
