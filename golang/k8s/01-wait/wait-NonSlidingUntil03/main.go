package main

import (
	"log"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
)

func main() {
	ch := make(chan struct{})
	i := 0

	go func() {
		log.Println("sleep 1s")
		time.Sleep(1 * time.Second)
		close(ch)
	}()

	wait.NonSlidingUntil(func() {
		time.Sleep(100 * time.Millisecond)
		i++
		log.Printf("%d test\n", i)
	}, 100 * time.Millisecond, ch)

	log.Printf("%d main exit\n", i)
}
