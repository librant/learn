package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

// sayHello 回调函数
func sayHello(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "hello k8s",
	})
}

func main() {
	log.SetFlags(log.Lshortfile)

	// 返回默认的路由引擎
	r := gin.Default()

	// 路由访问
	r.GET("/hello", sayHello)

	r.Run(":9090")
}
