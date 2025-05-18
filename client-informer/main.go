package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/workqueue"
)

func main() {
	// 获取 kubeconfig 文件路径
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

	// 创建 informer 工厂
	informerFactory := informers.NewSharedInformerFactory(clientset, time.Minute)

	// 创建 informer 对象
	podInformer := informerFactory.Core().V1().Pods()

	// 创建工作队列
	// 注意：这里使用的是 DefaultControllerRateLimiter()，这是一个默认的速率限制器，用于控制队列中的事件处理速度
	// 使用默认的速率限制器和配置创建工作队列
	// queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	rateLimite := workqueue.DefaultTypedControllerRateLimiter[any]()
	rateLimitingQueueConfig := workqueue.TypedRateLimitingQueueConfig[any]{}
	queue := workqueue.NewTypedRateLimitingQueueWithConfig(rateLimite, rateLimitingQueueConfig)

	// 定义处理新增、更新和删除事件的回调函数
	podHandler := cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err == nil {
				queue.Add(key)
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(newObj)
			if err == nil {
				queue.Add(key)
			}
		},
		DeleteFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err == nil {
				queue.Add(key)
			}
		},
	}

	// 将回调函数注册到 informer 上
	podInformer.Informer().AddEventHandler(podHandler)

	// 启动 informer
	stopCh := make(chan struct{})
	defer close(stopCh)
	informerFactory.Start(stopCh)

	// 等待 informer 同步完成
	if !cache.WaitForCacheSync(stopCh) {
		panic("同步 informer 缓存失败")
	}

	// 创建信号处理程序，用于捕捉 SIGTERM 和 SIGINT 信号
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGTERM, syscall.SIGINT)

	// 创建 worker 函数，用于处理队列中的事件
	processNextItem := func() {
		obj, shutdown := queue.Get()
		if shutdown {
			return
		}

		// 转换对象为 Pod
		key := obj.(string)
		podObj, exists, err := podInformer.Informer().GetIndexer().GetByKey(key)
		if err != nil {
			queue.Forget(obj)
			panic(fmt.Sprintf("获取 Pod 失败：%v", err))
		}

		if !exists {
			// 如果对象已经被删除，就把它从队列中移除
			queue.Forget(obj)
			return
		}

		// 在这里添加处理 Pod 的逻辑
		pod := podObj.(*v1.Pod)
		fmt.Printf("处理 Pod: namespace:%v,podName:%v\n", pod.Namespace, pod.Name)

		// 处理完事件后，把它从队列中移除
		queue.Forget(obj)
		return
	}

	// 启动 worker
	go wait.Until(processNextItem, time.Second, stopCh)

	// 等待信号
	<-signalCh
}
