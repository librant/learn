package main

import (
	"context"
	"flag"
	"log"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	log.SetFlags(log.Llongfile)
	log.Println("dynamic-client demon")

	// 通过参数传入 config 路径
	kubeconfig := flag.String("kubeconfig", "./.kube/kubeconfig",
		"Path to a kube config")
	flag.Parse()

	// 1 加载配置文件，生成 config 对象
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		log.Panicln(err)
	}

	// 2 生成 dynamic client 实例
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		log.Panicln(err)
	}

	// 3 设置请求的 gvr
	gvr := schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "pods",
	}

	// 4 发送请求， 且得到返回结果
	unStructData, err := dynamicClient.
		Resource(gvr).
		Namespace("default").
		List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Panicln(err)
	}

	// 5 转换为结构化数据
	podList := &corev1.PodList{}
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(
		unStructData.UnstructuredContent(), podList); err != nil {
		log.Panicln(err)
	}

	for _, item := range podList.Items {
		log.Printf("namespace: %s name: %s\n", item.Namespace, item.Name)
	}
}
