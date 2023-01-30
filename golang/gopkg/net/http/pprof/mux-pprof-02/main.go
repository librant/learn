package main

import (
	"log"
	"net/http"
	"net/http/pprof"
)

func InstallHandlerForPProf(mux *http.ServeMux) {
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/metrics", nil)
	InstallHandlerForPProf(mux)

	s := http.Server{
		Addr:    ":12345",
		Handler: mux,
	}
	log.Fatal(s.ListenAndServe())
}
