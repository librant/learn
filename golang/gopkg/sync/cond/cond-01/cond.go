package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	fmt.Println("Cond test...")

	c := sync.NewCond(&sync.Mutex{})
	ready := 0

	for i := 0; i < 10; i++ {
		go func(i int) {
			time.Sleep(time.Second * time.Duration(rand.Int63n(10)))
			// 加锁更改等待条件
			c.L.Lock()
			ready++
			c.L.Unlock()

			fmt.Printf("运动员%d已准备就绪\n",i)
			// 广播唤醒等待者，这里可以使用Broadcast和Signal
			c.Signal()
		}(i)
	}
	// 当修改条件或者 wait() 时，必须加锁，保护 condition
	for ready != 10 {
		c.L.Lock()
		c.Wait()
		c.L.Unlock()
		fmt.Println("裁判员被唤醒一次")
	}

	fmt.Println("所有运动员都准备就绪，比赛开始。。。")
}
