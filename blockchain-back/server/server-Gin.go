package server

import (
	"context"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime/trace"
	"syscall"
	"time"

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

	// Svg
	r.GET("/svg/:Cname", func(c *gin.Context) {
		trace.Log(c.Request.Context(), "Handler", "/svg/:Cname - Start")
		Svg := controllers.TakeBlock(c, c.Param("Cname"))
		c.Header("Content-Type", "svg+xml")
		c.Data(http.StatusOK, "svg+xml", Svg)
		trace.Log(c.Request.Context(), "Handler", "/svg/:Cname - End")
	})

	// Pdf
	r.GET("/pdf/:Cname", func(c *gin.Context) {
		trace.Log(c.Request.Context(), "Handler", "/pdf/:Cname - Start")
		PdfData := controllers.TakeBlock(c, c.Param("Cname"))
		c.Header("Content-Type", "application/pdf")
		c.Header("Content-Disposition", "inline; filename=example.pdf")
		c.Data(http.StatusOK, "application/pdf", PdfData)
		trace.Log(c.Request.Context(), "Handler", "/pdf/:Cname - End")
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

	srv := &http.Server{Addr: ":8080", Handler: r}
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Println("Starting server at :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
