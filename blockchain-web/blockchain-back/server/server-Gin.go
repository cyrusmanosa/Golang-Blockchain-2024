package server

import (
	"blockchain-back/controllers"
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Gin() {
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

	r.POST("/take", controllers.AddBlockForGin)
	r.POST("/Check/:name", controllers.AddBlockForGinConfirm)

	log.Println("Starting server at :8080")
	r.Run(":8080")
}
