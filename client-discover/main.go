package main

import (
	"fmt"
	"path/filepath"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	// 使用 kubeconfig 文件启动客户端
	dirPath := filepath.Join(homedir.HomeDir(), ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", dirPath)
	if err != nil {
		panic(err)
	}

	// 生成 discovery 客户端
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		panic(err)
	}
	// 获取 API 资源列表
	_, apiResourceList, err := discoveryClient.ServerGroupsAndResources()
	if err != nil {
		panic(err)
	}
	for _, v := range apiResourceList {
		// 解析 group/version
		gv, err := schema.ParseGroupVersion(v.GroupVersion)
		if err != nil {
			panic(err)
		}
		// 遍历资源列表，打印资源信息
		for _, resource := range v.APIResources {
			fmt.Printf("name:%v group:%v version:%v\n", resource.Name, gv.Group, gv.Version)

		}
	}
}
