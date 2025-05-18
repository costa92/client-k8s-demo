# Kubernetes 客户端示例项目

该项目包含多个使用 `client-go` 库与 Kubernetes API 交互的示例。

![示例图片](https://img2023.cnblogs.com/blog/2344773/202302/2344773-20230227230647754-846028579.png "示例图片说明")

## 项目结构

- `client-demo-one`: 动态客户端示例，列出 `default` 命名空间下的所有 Pod。
- `client-demo`: 其他客户端示例。
- `client-discover`: 使用 `discovery` 客户端获取 API 资源列表。
- `client-indexer`: 项目依赖管理。
- `client-informer`: 使用 `informer` 监听 Pod 事件。
- `client-lister`: 其他客户端示例。

## 依赖版本

- Go: 1.24.1
- `k8s.io` 相关库: v0.33.1

## 运行说明

每个子目录下的 `main.go` 文件均可独立运行，需要确保 `KUBECONFIG` 环境变量已正确设置，或者 `.kube/config` 文件存在。

例如，运行 `client-informer` 示例：

```sh
cd client-informer
go run main.go
```

## 参考

[client-go功能介绍](https://www.cnblogs.com/daemon365/p/17162339.html) 