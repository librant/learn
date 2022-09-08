package main

import (
	"log"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// 因为 informer 是一个持久运行的 goroutine，channel 作用：进程退出前通知 informer 退出
	stopChan := make(chan struct{})
	defer close(stopChan)

	// 获取 kube-config 文件
	cfg, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeDir)
	if err != nil {
		log.Fatalln(err)
	}
	// 生成 k8s clientSet
	clientSet, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		log.Fatalln(err)
	}

	// 第一步：创建 sharedInformer 象，第二个参数为重新同步数据的间隔时间
	sharedInformers := informers.NewSharedInformerFactory(clientSet, time.Minute)
	// 第二步：每个资源都有 informer 对象，这里获取 deployment 资源的 informer 对象
	deploymentInformer := sharedInformers.Apps().V1().Deployments().Informer()
	// 第三步：添加自定义回调函数
	deploymentInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		// 添加资源的回调函数，返回的是接口类型，需要强制转换为真正的类型
		AddFunc: func(obj interface{}) {
			mObj := obj.(v1.Object)
			log.Printf("New deployment added: %s", mObj.GetName())
		},
		// 更新资源的回调函数
		UpdateFunc: func(oldObj, newObj interface{}) {
			oObj := oldObj.(v1.Object)
			nObj := newObj.(v1.Object)
			log.Printf("%s deployment updated to %s", oObj.GetName(), nObj.GetName())
		},
		// 删除资源的回调函数
		DeleteFunc: func(obj interface{}) {
			mObj := obj.(v1.Object)
			log.Printf("deployment deleted from store: %s", mObj.GetName())
		},
	})
	// 第四步：开始运行informer对象
	deploymentInformer.Run(stopChan)
}
