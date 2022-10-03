package main

import (
	"flag"
	"log"

	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	log.SetFlags(log.Lshortfile)
	log.Printf("pod-informer demo")

	stopChan := make(chan struct{})
	defer close(stopChan)

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

	// 3 初始化 NewSharedInformerFactory
	sharedInformerFactory := informers.NewSharedInformerFactory(clientset, 0)

	// 生成 pod informers
	podInformer := sharedInformerFactory.Core().V1().Pods()
	// 生成 Indexer
	podIndexer := podInformer.Lister()

	// 等待数据同步完成
	sharedInformerFactory.WaitForCacheSync(stopChan)

	// 启动 informer
	sharedInformerFactory.Start(stopChan)

	// 利用 indexer 获取数据
	pods, err := podIndexer.List(labels.Everything())
	if err != nil {
		log.Panicln(err)
	}

	for _, item := range pods.Items {
		log.Printf("namespace: %s name: %s\n", item.Namespace, item.Name)
	}
}
