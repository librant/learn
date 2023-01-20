package main

import (
	"log"
	"strconv"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
)

// CusPod
type CusPod struct {
	ID   int
	Name string
}

func main() {
	podKillCh := make(chan *CusPod, 50)

	go func() {
		i := 0
		for {
			// 每 2s 向通道中发送当前 pod 的信息
			time.Sleep(2 * time.Second)
			podKillCh <- &CusPod{
				ID: i,
				Name: strconv.Itoa(i),
			}
			i++
		}
	}()

	wait.Until(func() {
		for pod := range podKillCh {
			log.Printf("%+v\n", pod)
		}
	}, 1*time.Second, wait.NeverStop)

	log.Println("main exit")
}
