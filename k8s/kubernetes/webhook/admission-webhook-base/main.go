package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"k8s.io/klog"

	"github.com/librant/learn/k8s/kubernetes/webhook/admission-webhook-base/webhook"
)

func main() {
	var parameters webhook.WhSvrParameters

	// get command line parameters
	flag.IntVar(&parameters.Port, "port", 443, "Webhook server port.")
	flag.StringVar(&parameters.CertFile, "tlsCertFile", "/etc/webhook/certs/cert.pem", "File containing the x509 Certificate for HTTPS.")
	flag.StringVar(&parameters.KeyFile, "tlsKeyFile", "/etc/webhook/certs/key.pem", "File containing the x509 private key to --tlsCertFile.")
	flag.Parse()

	whsvr, err := webhook.New(parameters)
	if err != nil {
		klog.Fatalln(err)
	}

	// define http server and server handler
	mux := http.NewServeMux()
	mux.HandleFunc("/mutate", whsvr.Serve)
	mux.HandleFunc("/validate", whsvr.Serve)
	whsvr.Svr.Handler = mux

	// start webhook server in new routine
	go func() {
		if err := whsvr.Svr.ListenAndServeTLS("", ""); err != nil {
			klog.Errorf("Failed to listen and serve webhook server: %v", err)
		}
	}()

	klog.Infof("Server started")

	// listening OS shutdown singal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	klog.Infof("Got OS shutdown signal, shutting down webhook server gracefully...")
	whsvr.Svr.Shutdown(context.Background())
}
