package main

import (
	"fmt"
	"log"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
)

func main() {
	log.SetFlags(log.Lshortfile)
	log.Printf("gin upload file demo")

	r := gin.Default()
	r.LoadHTMLFiles("./index.html")

	r.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.POST("/upload", func(c *gin.Context) {
		f, err := c.FormFile("f1")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		dst := path.Join("./", f.Filename)
		c.SaveUploadedFile(f, dst)
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("%s Uploaded", f.Filename),
		})
	})

	if err := r.Run(":9090"); err != nil {
		log.Panicln(err)
	}
}
