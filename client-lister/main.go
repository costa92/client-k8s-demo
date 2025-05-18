package main

import (
	"fmt"
	"os"
	"time"

	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	kubeconfigPath := os.Getenv("KUBECONFIG")
	if kubeconfigPath == "" {
		kubeconfigPath = os.Getenv("HOME") + "/.kube/config"
	}

	// 使用 kubeconfig 文件创建 kubernetes 客户端
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// 创建 Informer
	factory := informers.NewSharedInformerFactory(clientset, time.Minute)
	podInformer := factory.Core().V1().Pods()

	// 创建 Lister
	listers := podInformer.Lister()

	// 等待 Informer 同步完成
	stopCh := make(chan struct{})
	defer close(stopCh)

	factory.Start(stopCh)
	cache.WaitForCacheSync(stopCh, podInformer.Informer().HasSynced)

	// 获取 namespace 为 default 的 Pod 列表
	pods, err := listers.Pods("default").List(labels.Everything())
	if err != nil {
		panic(err)
	}
	for _, pod := range pods {
		fmt.Printf("Pod: %s\n", pod.Name)
	}
}
