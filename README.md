# neverdown

## 说明

一个 K8S 管理平台

## 架构图

## 模块说明

| 模块           | 说明                      |
| :------------- | :------------------------ |
| discovery      | 服务发现                  |
| authx          | 权限管理                  |
| authx_frontend | 权限管理, 前端            |
| jingx          | 镜像管理                  |
| jingx_frontend | 镜像管理, 前端            |
| openx          | 负责 K8S 集群的管理       |
| openx_frontend | 负责 K8S 集群的管理, 前端 |

## 安装部署

### 本地测试

本地测试环境通过 kind 来创建一个 K8S 集群 , 需要安装 metrics-server 用来获取资源使用情况

```sh
# 部署一个 K8S 集群
kind create cluster --config=./config/kind.yml

# 部署 metrics-server
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml

# 等待 metrics-server 启动 可能会面临如下的错误

Events:
  Type     Reason     Age                   From               Message
  ----     ------     ----                  ----               -------
 ...
 ...
 ...
  Warning  Unhealthy  59s (x31 over 5m29s)  kubelet            Readiness probe failed: HTTP probe failed with statuscode: 500

# 修改 metrics-server 的配置 进行如下修改操作

kubectl -n kube-system edit deploy metrics-server
spec:
      containers:
      - args:
        - --cert-dir=/tmp
        - --secure-port=10250
        - --kubelet-preferred-address-types=InternalIP,ExternalIP,Hostname
        - --kubelet-use-node-status-port
        - --metric-resolution=15s
        - --kubelet-insecure-tls

# 等待 metrics-server 启动正常
kubectl get pods -n kube-system | grep metrics

```

在 K8S 环境部署好之后, 需要外部启动 etcd 服务

然后启动 discovery, authx-apiserver, authx-frontend, openx-apiserver, openx-frontend 来体验

```sh
# 直接通过 docker 启动一个 etcd 服务
docker run -d --name etcd-server -p 2379:2379 -p 2380:2380 \
--env ALLOW_NONE_AUTHENTICATION=yes bitnami/etcd

# 部署 CRD
kubectl apply -f ./config/crd/

# 启动 discovery
make run-discovery-local

# 启动 authx apiserver
make run-authx-local

# 启动 authx frontend
make run-authx-frontend-local

# 启动 openx-apiserver
make run-openx-local

# 启动 openx-frontend
make run-openx-frontend-local
```

### 线上部署

线上部署, 我们这边是使用的阿里云的 ACK
