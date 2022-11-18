package main

import (
	"io"
	"log"
	"net/http"

	"github.com/emicklei/go-restful"
)

func main() {
	ws := new(restful.WebService)

	ws.Route(ws.GET("/hello").To(hello))
	restful.Add(ws)
	log.Fatalln(http.ListenAndServe(":12345", nil))
}

func hello(req *restful.Request, resp *restful.Response) {
	io.WriteString(resp, "world\n")
}
