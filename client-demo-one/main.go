package main

import (
	"context"
	"fmt"
	"path/filepath"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	ctx := context.Background()

	// 加载kubeconfig配置文件，生成config对象
	dirPath := homedir.HomeDir()
	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(dirPath, ".kube", "config"))
	if err != nil {
		panic(err)
	}

	// dynamicClient
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// 定义组版本资源
	gvr := schema.GroupVersionResource{Version: "v1", Resource: "pods"}
	unStructObj, err := dynamicClient.Resource(gvr).Namespace("default").List(ctx, metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	podList := &apiv1.PodList{}

	if err = runtime.DefaultUnstructuredConverter.FromUnstructured(unStructObj.UnstructuredContent(), podList); err != nil {
		panic(err)
	}

	for _, v := range podList.Items {
		fmt.Printf("Namespace:%v podName:%v status:%v \n", v.Namespace, v.Name, v.Status.Phase)
	}
}
