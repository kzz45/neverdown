# 使用说明

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

```sh
# 直接通过 docker 启动一个 etcd 服务
docker run -d --name etcd-server -p 2379:2379 -p 2380:2380 \
--env ALLOW_NONE_AUTHENTICATION=yes bitnami/etcd

# 启动 discovery
make run-discovery-local

# 启动 authx apiserver
make run-authx-local

# 启动 authx frontend
make run-authx-frontend-local

# 部署 CRD
kubectl apply -f config/crd/

# 启动 openx-apiserver
make run-openx-local

# 启动 openx-frontend
make run-openx-frontend-local
```

### 线上部署

线上部署, 我们这边是使用的阿里云的 ACK

```sh
kubectl apply -f config/crd/
kubectl apply -f config/deploy/neverdown-deploy.yaml
```
