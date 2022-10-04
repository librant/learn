package main

import (
	"log"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
)

func main() {
	ch := make(chan struct{})
	go func() {
		log.Println("sleep 1s")
		time.Sleep(1 * time.Second)
		close(ch)
	}()
	wait.Until(func() {
		time.Sleep(100 * time.Millisecond)
		log.Println("test")
	}, 100*time.Millisecond, ch)
	log.Println("main exit")
}
