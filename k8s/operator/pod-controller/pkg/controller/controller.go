package controller

import (
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	informer "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	lister "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
)

type controller struct {
	client    kubernetes.Interface
	podLister lister.PodLister
	podSynced cache.InformerSynced
	queue     workqueue.RateLimitingInterface
	recorder  record.EventRecorder
}

// New 实例化
func New(clientset kubernetes.Interface, podInformer informer.PodInformer) *controller {
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(klog.Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{
		Interface: clientset.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme,
		corev1.EventSource{Component: controllerAgentName})

	c := controller{
		client:    clientset,
		podLister: podInformer.Lister(),
		podSynced: podInformer.Informer().HasSynced,
		recorder:  recorder,
		queue:     workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "pod"),
	}

	podInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    c.addPod,
		UpdateFunc: c.uptPod,
		DeleteFunc: c.delPod,
	})

	return &c
}

// Run 启动 controller
func (c *controller) Run(threadiness int, stopCh <-chan struct{}) error {
	defer runtime.HandleCrash()
	defer c.queue.ShutDown()

	// 开始处理业务逻辑，开始同步数据
	klog.Infof("begin pod sync...")
	if ok := cache.WaitForCacheSync(stopCh, c.podSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	klog.Infof("begin worker start...")
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Minute, stopCh)
	}

	klog.Info("Started workers")
	<-stopCh
	klog.Info("Stop workers")
	return nil
}

func (c *controller) runWorker() {
	for c.processNextWorkItem() {
	}
}

func (c *controller) processNextWorkItem() bool {
	item, shutdown := c.queue.Get()
	if shutdown {
		return false
	}

	if err := func(obj interface{}) error {
		defer c.queue.Done(obj)

		key, ok := obj.(string)
		if !ok {
			c.queue.Forget(obj)
			runtime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}
		// 在 syncHandler 中处理业务
		if err := c.syncHandler(key); err != nil {
			return fmt.Errorf("error syncing '%s': %s", key, err.Error())
		}
		c.queue.Forget(obj)
		klog.Infof("Successfully synced '%s'", key)
		return nil
	}(item); err != nil {
		runtime.HandleError(err)
	}
	return true
}

func (c *controller) syncHandler(key string) error {
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		runtime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	pod, err := c.podLister.Pods(namespace).Get(name)
	if err != nil {
		if errors.IsNotFound(err) {
			// TODO: 这里 pod 被删除了，这里可以处理实际逻辑

			klog.Infof("pod name: %s namespace: %s has been delete", name, namespace)
			return nil
		}
		runtime.HandleError(fmt.Errorf("failed to list pod by: %s/%s", namespace, name))
		return err
	}

	// 这是实际获取集群中 pod 的状态
	klog.Infof("pod yaml: %s/%s", pod.Namespace, pod.Name)

	// TODO：这里可以校验当前 pod 的状态

	c.recorder.Event(pod, corev1.EventTypeNormal, successSynced, messageResourceSynced)
	return nil
}

func (c *controller) addPod(obj interface{}) {
	pod := obj.(*corev1.Pod)
	klog.Infof("add pod name: %s namespace: %s", pod.Name, pod.Namespace)
	c.enqueuePod(obj)
}

func (c *controller) uptPod(oldObj, newObj interface{}) {
	oldPod := oldObj.(*corev1.Pod)
	newPod := newObj.(*corev1.Pod)

	if oldPod.ResourceVersion == newPod.ResourceVersion {
		// 版本一致，就表示没有实际更新的操作，立即返回
		return
	}

	c.enqueuePod(newPod)
}

func (c *controller) delPod(obj interface{}) {
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil {
		runtime.HandleError(err)
		return
	}
	// 删除成功
	c.queue.AddRateLimited(key)
}

// enqueuePod 数据先放入缓存，再入队列
func (c *controller) enqueuePod(obj interface{}) {
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err != nil {
		runtime.HandleError(err)
		return
	}
	// 将 key 放入队列
	c.queue.AddRateLimited(key)
}
