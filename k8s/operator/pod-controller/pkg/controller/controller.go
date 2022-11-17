package controller

import (
	"log"
	"time"

	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	listerv1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
)

type controller struct {
	client    kubernetes.Interface
	podLister listerv1.PodLister
	queue     workqueue.RateLimitingInterface
}

// New 实例化
func New(kubeConfig string) *controller {
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		log.Panicln(err)
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Panicln(err)
	}
	return &controller{
		client: clientSet,
		queue:  workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter()),
	}
}

// Run 启动 controller
func (c *controller) Run(stopCh chan struct{}) {
	defer runtime.HandleCrash()
	defer c.queue.ShutDown()

	factory := informers.NewSharedInformerFactory(c.client, 0)
	go factory.Start(stopCh)

	podInformer := factory.Core().V1().Pods()
	c.podLister = podInformer.Lister()

	podInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    c.addPod,
		UpdateFunc: c.uptPod,
		DeleteFunc: c.delPod,
	})

	if ok := cache.WaitForCacheSync(stopCh, podInformer.Informer().HasSynced); !ok {
		log.Panicln("HasSynced failed")
	}

	for i := 0; i < workNum; i++ {
		go wait.Until(c.runWorker, time.Minute, stopCh)
	}

	klog.Info("Started workers")
	<-stopCh
}

func (c *controller) runWorker() {

}

func (c *controller) addPod(obj interface{}) {

}

func (c *controller) uptPod(oldObj, newObj interface{}) {

}

func (c *controller) delPod(obj interface{}) {

}
