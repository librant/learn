package main

import (
	"log"

	"github.com/librant/learn/k8s/operator/service-ingress-controller/cmd"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ltime)
	log.Printf("service-ingress-controller start...")

	cmd.Execute()
}
