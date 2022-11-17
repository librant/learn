package main

import (
	"log"

	"github.com/librant/learn/k8s/operator/pod-controller/cmd"
)

var (
	Author = "librant"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ltime)

	log.Printf("pod-controller author: %s", Author)

	cmd.Execute()
}
