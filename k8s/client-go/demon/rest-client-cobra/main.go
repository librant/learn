package main

import (
	"log"

	"github.com/librant/learn/k8s/client-go/demon/rest-client-cobra/cmd"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ltime)
	log.Printf("hello world")

	cmd.Execute()
}
