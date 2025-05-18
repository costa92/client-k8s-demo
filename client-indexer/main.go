package main

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
)

const (
	NamespaceIndexName = "namespace" // 定义一个索引器名称，用于按照命名空间检索Pod
	NodeNameIndexName  = "nodeName"  // 定义一个索引器名称，用于按照节点名称检索Pod
)

// NamespaceIndexFunc是一个函数，用于从对象中提取命名空间作为索引键
func NamespaceIndexFunc(obj interface{}) ([]string, error) {
	m, err := meta.Accessor(obj) // 从对象中提取元数据
	if err != nil {
		return []string{""}, fmt.Errorf("object has no meta: %v", err)
	}
	return []string{m.GetNamespace()}, nil // 返回对象的命名空间作为索引键
}

// NodeNameIndexFunc是一个函数，用于从对象中提取节点名称作为索引键
func NodeNameIndexFunc(obj interface{}) ([]string, error) {
	pod, ok := obj.(*v1.Pod) // 将对象转换为Pod类型
	if !ok {
		return []string{""}, fmt.Errorf("object is not a Pod")
	}

	return []string{pod.Spec.NodeName}, nil // 返回Pod的节点名称作为索引键
}

func main() {
	index := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{
		NamespaceIndexName: NamespaceIndexFunc, // 使用NamespaceIndexFunc作为命名空间索引器
		NodeNameIndexName:  NodeNameIndexFunc,  // 使用NodeNameIndexFunc作为节点名称索引器
	}) // 创建一个新的索引器，指定主键函数和辅助键函数

	pod1 := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "index-pod-1",
			Namespace: "default",
		},
		Spec: v1.PodSpec{
			NodeName: "orbstack",
		},
	}
	pod2 := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "index-pod-2",
			Namespace: "default",
		},
		Spec: v1.PodSpec{NodeName: "orbstack"},
	} //

	index.Add(pod1) // 将pod1添加到索引器中
	index.Add(pod2) // 将pod2添加到索引器中

	pods, err := index.ByIndex(NamespaceIndexName, "default") // 按照命名空间为default检索Pod列表
	if err != nil {
		panic(err)
	}
	for _, pod := range pods {
		fmt.Println(pod.(*v1.Pod).Name)
	} // 遍历并打印检索到的Pod名称
	fmt.Println("*****************")
	pods, err = index.ByIndex(NodeNameIndexName, "orbstack") // 按照节点名称为node2检索Pod列表
	if err != nil {
		panic(err)
	}
	for _, pod := range pods {
		fmt.Println(pod.(*v1.Pod).Name)
	} // 遍历并打印
}
