package main

import (
	"fmt"
	"log"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
)

// NamespaceIndexFunc namespace 索引器函数
func NamespaceIndexFunc(obj interface{}) ([]string, error) {
	pod, ok := obj.(*corev1.Pod)
	if !ok {
		return nil, fmt.Errorf("kind error")
	}
	return []string{pod.Namespace}, nil
}

// NodeIndexFunc node 索引器函数
func NodeIndexFunc(obj interface{}) ([]string, error) {
	pod, ok := obj.(*corev1.Pod)
	if !ok {
		return nil, fmt.Errorf("kind error")
	}
	return []string{pod.Spec.NodeName}, nil
}

func main() {
	// 生成 index 实例
	index := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{
		"namespace": NamespaceIndexFunc,
		"node":      NodeIndexFunc,
	})

	// 构造数据
	pods := getPods(3, 0)

	// 将对象加入到 Indexer 中
	for _, pod := range pods {
		if err := index.Add(pod); err != nil {
			log.Printf("index add %v failed: %v", pod.Name, err)
		}
	}

	// 通过索引器函数查询数据
	podList, err := index.ByIndex("namespace", "default")
	if err != nil {
		log.Panicln(err)
	}

	for _, item := range podList {
		pod := item.(*corev1.Pod)
		log.Printf("name: %s\n", pod.Name)
	}

	nodeList, err := index.ByIndex("node", "node-1")
	if err != nil {
		log.Panicln(err)
	}

	for _, item := range nodeList {
		pod := item.(*corev1.Pod)
		log.Printf("node: %s\n", pod.Spec.NodeName)
	}
}

func getPods(num, start int) []*corev1.Pod {
	pods := make([]*corev1.Pod, num)
	for i := 0; i < num; i++ {
		pods[i] = &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("pod-index-%d", i+start),
				Namespace: getNamespaceMap(i),
			},
			Spec: corev1.PodSpec{
				NodeName: fmt.Sprintf("node-%d", i+start),
			},
		}
	}
	return pods
}

var namespaceMap = map[int]string{
	0: "default",
	1: "kube-system",
	2: "kube-public",
}

func getNamespaceMap(index int) string {
	return namespaceMap[index]
}
