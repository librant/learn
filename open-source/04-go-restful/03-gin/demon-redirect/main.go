package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	log.SetFlags(log.Lshortfile)

	r := gin.Default()

	r.GET("/hello", func(c *gin.Context) {
		// 指定到具体的 location
		c.Redirect(http.StatusMovedPermanently, "https://www.baidu.com")
	})

	r.GET("/a", func(c *gin.Context) {
		c.Request.URL.Path = "/b"
		r.HandleContext(c)
	})

	r.GET("/b", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "b",
		})
	})

	if err := r.Run(":9090"); err != nil {
		log.Panicln(err)
	}
}
