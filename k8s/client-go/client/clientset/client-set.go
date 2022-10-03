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
	log.Println("client-set demon")

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

	// 3 查询 default 命名空间下的 pods 信息
	pods, err := clientset.
		CoreV1().
		Pods("default").
		List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Panicln(err)
	}

	for _, item := range pods.Items {
		log.Printf("namespace: %s name: %s\n", item.Namespace, item.Name)
	}
}
