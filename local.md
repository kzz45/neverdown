# 本地测试

- 本地测试环境通过 kind 来创建一个 K8S 集群

![install k8s](./doc/install_k8s_with_kind.png)

- 需要安装 metrics-server 用来获取资源使用情况

![install metrics-server](./doc/install_metrics_server.png)

- 在 K8S 环境部署好之后, 需要外部启动 etcd 服务

![install etcd](./doc/install_etcd_out_cluster.png)

- 安装 CRD

![install crd](./doc/install_crd.png)

- 安装 RBAC Role
  
![install rbac](./doc/install_rbac.png)

- 启动 disovery 服务

![run-discovery-local](./doc/run-discovery-local.png)

- 启动 authx-apiserver

![run-authx-local](./doc/run-authx-local.png)

- 启动 jingx-apiserver

![run-jingx-local](./doc/run-jingx-local.png)

- 启动 openx-apiserver

![run-openx-local](./doc/run-openx-local.png)

- 启动 authx-frontend 并复制 openx-apiserver 的密码

```sh
make run-authx-frontend-local
```

![run-authx-frontend](./doc/run-authx-frontend.png)

- 启动 openx-frontend 并登录

```sh
make run-openx-frontend-local
```

![run-openx-frontend-local](./doc/run-openx-frontend-local.png)

- 测试创建 nginx 服务

![nginx-1](./doc/nginx-1.png)

![nginx-2](./doc/nginx-2.png)

![nginx-3](./doc/nginx-3.png)

![nginx-4](./doc/nginx-4.png)

![nginx-5](./doc/nginx-5.png)

![nginx-6](./doc/nginx-6.png)

![nginx-7](./doc/nginx-7.png)

![nginx-8](./doc/nginx-8.png)
