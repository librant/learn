package pkg

import (
	v13 "k8s.io/client-go/informers/core/v1"
	v14 "k8s.io/client-go/informers/networking/v1"
	"k8s.io/client-go/kubernetes"
	v12 "k8s.io/client-go/listers/core/v1"
	v11 "k8s.io/client-go/listers/networking/v1"
	"k8s.io/client-go/tools/cache"
)

type controller struct {
	client        kubernetes.Interface
	serviceLister v12.ServiceLister
	ingressLister v11.IngressLister
}

func (c *controller) CreateService(obj interface{}) {

}

func (c *controller) UpdateService(oldObj interface{}, newObj interface{}) {

}

func (c *controller) DeleteIngress(obj interface{}) {

}

func (c *controller) Run(stopCh chan struct{}) {
	<- stopCh
}

func New(client kubernetes.Interface,
	serviceInformer v13.ServiceInformer,
	ingressInformer v14.IngressInformer) *controller {

	c := controller{
		client:        client,
		serviceLister: serviceInformer.Lister(),
		ingressLister: ingressInformer.Lister(),
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
