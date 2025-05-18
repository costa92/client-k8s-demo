# Kubernetes 客户端示例项目

该项目包含多个使用 `client-go` 库与 Kubernetes API 交互的示例，涵盖了动态客户端、资源发现、Informer 机制、Indexer、Lister 及自定义控制器等核心用法。

![client-go 架构图](https://img2023.cnblogs.com/blog/2344773/202302/2344773-20230227230647754-846028579.png "client-go informer/controller 架构流程")

## client-go Informer/Controller 工作机制简介

如上图所示，`client-go` 的核心机制包括以下几个部分：

1. **Reflector**：通过 List & Watch 机制从 Kubernetes API 服务器同步资源对象，推送到 Delta Fifo 队列。
2. **Informer**：从队列中 Pop 出对象，分发给 Indexer，并触发事件处理器（Event Handlers）。
3. **Indexer**：负责对象的本地缓存和索引，支持高效查询。
4. **Resource Event Handlers**：监听资源的增删改事件，将对象的 Key 入队到 Workqueue。
5. **Workqueue & Controller**：自定义控制器从队列中取出 Key，调用业务逻辑处理（Handle Object）。

该机制实现了高效、可靠的资源监听与处理，是自定义 Kubernetes 控制器的基础。

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

## 参考资料

- [client-go功能介绍](https://www.cnblogs.com/daemon365/p/17162339.html)
- [client-go 官方文档](https://github.com/kubernetes/client-go)
- [Kubernetes 官方文档](https://kubernetes.io/zh/docs/home/)

---

**补充说明**：  
本项目各模块均遵循 Go 语言最佳实践，代码结构清晰，便于扩展和集成。建议结合上方架构图理解 informer/controller 的事件流转过程，有助于开发高效、健壮的自定义控制器。

如需进一步集成 Python、Next.js 或 LangChain 等技术栈，可参考相关官方文档，结合微服务架构和现代云原生开发最佳实践进行扩展。 