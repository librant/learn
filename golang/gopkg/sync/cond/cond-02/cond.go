package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	fmt.Println("Cond test...")

	shared := false
	c := sync.NewCond(&sync.Mutex{})

	wg := sync.WaitGroup{}
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(i int) {
			// 等待 shared 为 true
			c.L.Lock()
			for !shared {
				fmt.Printf("time: %v gorouting[%d] wait\n", time.Now(), i)
				c.Wait()
			}
			c.L.Unlock()
			wg.Done()
		}(i)
	}

	// 更改 shared 为 true
	time.Sleep(2 * time.Second)

	c.L.Lock()
	fmt.Printf("time: %v, main gorouting ready\n", time.Now().Format(time.RFC3339))
	shared = true
	c.L.Unlock()
	c.Broadcast()
	fmt.Printf("time: %v, main gorouting broadcast\n", time.Now().Format(time.RFC3339))
	wg.Wait()
}
