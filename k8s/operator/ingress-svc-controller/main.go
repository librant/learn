package main

import (
	"k8s.io/klog"

	"github.com/librant/learn/k8s/operator/ingress-svc-controller/cmd"
)

var (
	Author = "librant"
)

func main() {
	klog.Infof("controller-controller author: %s", Author)

	cmd.Execute()
}
