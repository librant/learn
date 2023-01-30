package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/librant/learn/golang/gopkg/net/http/pprof/debug-pprof-01/data"
)

// http://127.0.0.1:12345/debug/pprof/

func main() {
	go func() {
		for {
			log.Println(data.Add("https://github.com/librant"))
		}
	}()
	log.Fatal(http.ListenAndServe(":12345", nil))
}