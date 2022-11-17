package controller

import (
	"log"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type controller struct {
	clientSet kubernetes.Interface
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
		clientSet: clientSet,
	}
}
