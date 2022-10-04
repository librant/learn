package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	log.SetFlags(log.Lshortfile)
	log.Printf("gin router demo")

	r := gin.Default()

	// 路由
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "github.com/librant",
		})
	})
	r.Any("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"method": c.Request.Method,
		})
	})

	// 路由组
	userGroup := r.Group("/vedio")
	userGroup.GET("/index", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "index",
		})
	})

	if err := r.Run(":9090"); err != nil {
		log.Panicln(err)
	}
}