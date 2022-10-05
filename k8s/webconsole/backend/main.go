package main

import (
	"flag"
	"log"

	"github.com/gin-gonic/gin"
)

var kubeconfig *string

func init() {
	kubeconfig = flag.String("kubeconfig", "./.kube/kubeconfig",
		"Path to a kube config")
}

func main() {
	log.SetFlags(log.Lshortfile)
	flag.Parse()

	r := gin.Default()

	if err := r.Run(":9090"); err != nil {
		log.Fatalln(err)
	}
}
