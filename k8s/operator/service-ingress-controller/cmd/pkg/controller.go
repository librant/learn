package pkg

import (
	"context"
	"reflect"
	"time"

	v10 "k8s.io/api/core/v1"
	v15 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v16 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	v13 "k8s.io/client-go/informers/core/v1"
	v14 "k8s.io/client-go/informers/networking/v1"
	"k8s.io/client-go/kubernetes"
	v12 "k8s.io/client-go/listers/core/v1"
	v11 "k8s.io/client-go/listers/networking/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type controller struct {
	client        kubernetes.Interface
	serviceLister v12.ServiceLister
	ingressLister v11.IngressLister
	queue         workqueue.RateLimitingInterface
}

func (c *controller) CreateService(obj interface{}) {
	c.enqueue(obj)
}

func (c *controller) UpdateService(oldObj interface{}, newObj interface{}) {
	// TODO: 比较 annotation 是否一致
	if reflect.DeepEqual(oldObj, newObj) {
		return
	}
	c.enqueue(newObj)
}

func (c *controller) DeleteIngress(obj interface{}) {
	ingress := obj.(*v15.Ingress)
	ownerReference := v16.GetControllerOf(ingress)
	if ownerReference == nil {
		return
	}
	if ownerReference.Kind != "Service" {
		return
	}
	c.queue.Add(ingress.Namespace + "/" + ingress.Name)
}

func (c *controller) enqueue(obj interface{}) {
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err != nil {
		runtime.HandleError(err)
	}
	c.queue.Add(key)
}

func (c *controller) Run(stopCh chan struct{}) {
	for i := 0; i < workNum; i++ {
		go wait.Until(c.worker, time.Minute, stopCh)
	}

	<-stopCh
}

func (c *controller) worker() {
	for c.processNextItem() {

	}
}

func (c *controller) processNextItem() bool {
	item, st := c.queue.Get()
	if !st {
		return false
	}
	defer c.queue.Done(item)

	key := item.(string)
	if err := c.syncService(key); err != nil {
		c.handleError(key, err)
	}
	return true
}

func (c *controller) handleError(key string, err error) {
	if c.queue.NumRequeues(key) < maxRetry {
		c.queue.AddRateLimited(key)
	}
	runtime.HandleError(err)
	c.queue.Forget(key)
}

func (c *controller) syncService(key string) error {
	nsKey, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}

	service, err := c.serviceLister.Services(nsKey).Get(name)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}

	_, ok := service.GetAnnotations()["ingress/http"]
	ingress, err := c.ingressLister.Ingresses(nsKey).Get(name)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	if ok && errors.IsNotFound(err) {
		// 创建 ingress
		ig := c.constructIngress(service)
		if _, err := c.client.NetworkingV1().Ingresses(nsKey).Create(context.Background(),
			ig, v16.CreateOptions{}); err != nil {
			return err
		}
	} else if !ok && ingress != nil {
		if err := c.client.NetworkingV1().Ingresses(nsKey).Delete(context.Background(),
			name, v16.DeleteOptions{}); err != nil {
			return err
		}
	}
	return nil
}

func (c *controller) constructIngress(service *v10.Service) *v15.Ingress {
	pathType := v15.PathTypePrefix
	ownRef := v16.NewControllerRef(service, v16.SchemeGroupVersion.WithKind("Service"))
	ingress := v15.Ingress{
		ObjectMeta: v16.ObjectMeta{
			Name:      service.Name,
			Namespace: service.Namespace,
			OwnerReferences: []v16.OwnerReference{
				*ownRef,
			},
		},
		Spec: v15.IngressSpec{
			Rules: []v15.IngressRule{
				{
					Host: "librant.example.com",
					IngressRuleValue: v15.IngressRuleValue{
						HTTP: &v15.HTTPIngressRuleValue{
							Paths: []v15.HTTPIngressPath{
								{
									Path:     "/",
									PathType: &pathType,
									Backend: v15.IngressBackend{
										Service: &v15.IngressServiceBackend{
											Name: service.Name,
											Port: v15.ServiceBackendPort{
												Name:   "http",
												Number: 8080,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	return &ingress
}

func New(client kubernetes.Interface,
	serviceInformer v13.ServiceInformer,
	ingressInformer v14.IngressInformer) *controller {

	c := controller{
		client:        client,
		serviceLister: serviceInformer.Lister(),
		ingressLister: ingressInformer.Lister(),
		queue:         workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter()),
	}

	serviceInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    c.CreateService,
		UpdateFunc: c.UpdateService,
	})

	ingressInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		DeleteFunc: c.DeleteIngress,
	})

	return &c
}
