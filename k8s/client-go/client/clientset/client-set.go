package main

import (
	"context"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	log.Println("clientset demon")

	// 1 加载配置文件
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeDir)
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
