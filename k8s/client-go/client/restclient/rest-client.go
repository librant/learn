package main

import (
	"context"
	"log"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	log.Println("rest-client demo")

	// 1 指定 k8s 的配置文件
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeDir)
	if err != nil {
		log.Panicln(err)
	}

	// 2 配置 API 路径
	config.APIPath = "api"

	// 3 设置分组版本
	config.GroupVersion = &corev1.SchemeGroupVersion

	// 4 配置数据的编解码器
	config.NegotiatedSerializer = scheme.Codecs

	// 5 实例化 rest client
	restClient, err := rest.RESTClientFor(config)
	if err != nil {
		log.Panicln(err)
	}

	// 6 定义返回接收值
	result := &corev1.PodList{}
	if err := restClient.Get().
		Namespace("kube-system").
		Resource("pods").
		VersionedParams(&metav1.ListOptions{}, scheme.ParameterCodec).
		Do(context.Background()).
		Into(result); err != nil {
		log.Panicln(err)
	}

	for _, item := range result.Items {
		log.Printf("namespace: %s name: %s\n", item.Namespace, item.Name)
	}
}
