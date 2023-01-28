package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

func main() {
	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			u, _ := url.Parse("https://www.baidu.com/")
			req.URL = u
			req.Host = u.Host // 必须显示修改 Host，否则转发可能失败
		},
		ModifyResponse: func(resp *http.Response) error {
			log.Println("resp status:", resp.Status)
			log.Println("resp headers:")
			for hk, hv := range resp.Header {
				log.Println(hk, ":", strings.Join(hv, ","))
			}
			return nil
		},
		ErrorLog: log.New(os.Stdout, "ReverseProxy:", log.LstdFlags | log.Lshortfile),
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			if err != nil {
				log.Println("ErrorHandler catch err:", err)
				w.WriteHeader(http.StatusBadGateway)
				_, _ = fmt.Fprintf(w, err.Error())
			}
		},
	}

	http.Handle("/", proxy)

	if err := http.ListenAndServe(":12345", nil); err != nil {
		log.Fatal(err)
	}
}
