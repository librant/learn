package main

import (
	"k8s.io/klog"

	"github.com/librant/learn/librant/project/cluster-console/cmd"
)

var (
	Author = "librant"
)

func main() {
	klog.Infof("cluster-console author: %s", Author)

	cmd.Execute()
}
