package main

import (
	"context"
	"flag"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	log.SetFlags(log.Lshortfile)
	log.Printf("reflector watch demo")

	// 通过参数传入 config 路径
	kubeconfig := flag.String("kubeconfig", "./.kube/kubeconfig",
		"Path to a kube config")
	flag.Parse()

	// 1 加载配置文件
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		log.Panicln(err)
	}

	// 2 实例化 clientset 对象
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Panicln(err)
	}

	// 3 建立 watch 的链接
	w, err := clientset.AppsV1().Deployments("default").
		Watch(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Panicln(err)
	}

	log.Printf("watch deployment begin...")

	// 4 实时监听 watch 资源
	for {
		select {
		case e, _ := <-w.ResultChan():
			log.Printf("type: %v object: %v", e.Type, e.Object)
		}
	}
}
