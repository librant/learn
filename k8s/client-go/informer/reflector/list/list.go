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
	log.Printf("reflector list demo")

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

	// 3 list 资源列表
	deploymentList, err := clientset.AppsV1().Deployments("default").
		List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Panicln(err)
	}

	// 4 打印获取的资源信息
	for _, item := range deploymentList.Items {
		log.Printf("namespace: %s kind: %s, name: %s\n",
			item.Namespace, item.Kind, item.Name)
	}
}