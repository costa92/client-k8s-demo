package main

import (
	"context"
	"fmt"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(homedir.HomeDir(), ".kube", "config"))
	if err != nil {
		fmt.Printf("获取配置失败: %v\n", err)
		return
	}

	config.APIPath = "api"
	config.GroupVersion = &corev1.SchemeGroupVersion
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()

	restClient, err := rest.RESTClientFor(config)
	if err != nil {
		fmt.Printf("创建 REST 客户端失败: %v\n", err)
		return
	}

	rest := &corev1.PodList{}
	if err = restClient.Get().Namespace("default").Resource("pods").VersionedParams(&metav1.ListOptions{},
		scheme.ParameterCodec).Do(context.TODO()).Into(rest); err != nil {
		fmt.Printf("获取 Pod 列表失败: %v\n", err)
		return
	}
	for _, v := range rest.Items {
		fmt.Printf("NameSpace: %v  Name: %v  Status: %v  Generation: %v ResourceVersion %v \n", v.Namespace, v.Name, v.Status.Phase, v.Generation, v.ResourceVersion)

		// 原代码中 fmt.Printf 格式字符串有两个占位符，但只提供了一个参数，这里修正为正确输出 OwnerReferences 的 Kind 和 APIVersion
		// GetOwnerReferences 返回的是一个 OwnerReference 切片，需要遍历输出
		// fmt.Printf("OwnerReferences: %v \n", v.GetOwnerReferences())
		for _, ref := range v.GetOwnerReferences() {
			fmt.Printf("Kind: %v  APIVersion: %v \n", ref.Kind, ref.APIVersion)
		}
	}
}
