package main

import (
	"flag"
	"log"
	"time"

	"k8s.io/client-go/discovery/cached/disk"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	log.SetFlags(log.Lshortfile)
	log.Println("cache-discovery-client demo")

	// 通过参数传入 config 路径
	kubeconfig := flag.String("kubeconfig", "./.kube/kubeconfig",
		"Path to a kube config")
	flag.Parse()

	// 1 加载配置文件， 生成 config 对象
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		log.Panicln(err)
	}

	// 2 实例化客户端，本客户端负载将 GVR 数据，缓存本地文件中
	cacheDiscoveryClient, err := disk.NewCachedDiscoveryClientForConfig(config,
		"./cache/discovery", "./cache/http", time.Minute*60)
	if err != nil {
		log.Panicln(err)
	}

	// 3 发送请求，获取数据
	_, _, err = cacheDiscoveryClient.ServerGroupsAndResources()
	// 1) 先从缓存文件中查找 gvr 数据，有则直接返回，没有则需要调用 APIServer
	// 2) 调用 APIServer 获取 GVR 数据
	// 3) 将获取的 GVR 数据缓存到本地，然后返回给数据端
	if err != nil {
		log.Panicln(err)
	}

	// 这个输出同 discovery-client demo
}
