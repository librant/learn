package main

import (
	"log"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	log.SetFlags(log.Llongfile)
	log.Println("discovery-client demo")

	// 1 加载配置文件， 生成 config 对象
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeDir)
	if err != nil {
		log.Panicln(err)
	}

	// 2 实例化客户端
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		log.Panicln(err)
	}

	// 3 发送数据， 获取 gvr
	_, apiResources, err := discoveryClient.ServerGroupsAndResources()
	if err != nil {
		log.Panicln(err)
	}

	// 4 解析数据
	for _, list := range apiResources {
		gv, err := schema.ParseGroupVersion(list.GroupVersion)
		if err != nil {
			log.Panicln(err)
		}

		for _, resource := range list.APIResources {
			log.Printf("name: %v group: %v version: %v",
				resource.Name, gv.Group, gv.Version)
		}
	}
}
