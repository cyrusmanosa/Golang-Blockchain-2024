package server

import (
	"blockchain-back/controllers"
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

	r.GET("/svg/:Cname", func(c *gin.Context) {
		Cname := c.Param("Cname")
		Svg := controllers.TakeBlock(c, Cname)
		c.Header("Content-Type", "svg+xml")
		c.Data(http.StatusOK, "svg+xml", Svg)
	})

	r.PUT("/take", controllers.AddBlockForGin)
	r.PUT("/Check/:name", controllers.AddBlockForGinConfirm)

	fmt.Println("Starting server at :8080")
	log.Fatal(r.Run(":8080"))
}
