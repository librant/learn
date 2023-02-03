package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var once sync.Once

	// 多个线程只执行一次
	onceBody := func() {
		fmt.Println("Only once")
	}

	done := make(chan bool)
	for i := 0; i < 5; i++ {
		j := i
		go func(int) {
			once.Do(onceBody)
			fmt.Println(j)
			done <- true
		}(j)
	}
	<-done

	time.Sleep(2 * time.Second)
}
