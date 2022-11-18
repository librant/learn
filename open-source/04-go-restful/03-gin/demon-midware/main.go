package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// indexHandler index handler
func indexHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "index",
	})
}

// midWareHandler mid ware
func midWareHandler(c *gin.Context) {
	fmt.Println("mid handler...")
	// 执行的起始时间
	start := time.Now()
	c.Next()
	cost := time.Since(start)
	log.Printf("cost: %v\n", cost)
}

func main() {
	log.SetFlags(log.Lshortfile)

	r := gin.Default()

	r.GET("/index", midWareHandler, indexHandler)

	if err := r.Run(":9090"); err != nil {
		log.Panicln(err)
	}
}
