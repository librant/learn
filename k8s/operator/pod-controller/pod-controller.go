package main

import (
	"log"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ltime)

	log.Printf("pod-controller start...")
}
