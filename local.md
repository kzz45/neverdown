# 本地测试

## 集群创建

- 本地测试环境通过 kind 来创建一个 K8S 集群

![install k8s](./doc/install_k8s_with_kind.png)

- 需要安装 metrics-server 用来获取资源使用情况

![install metrics-server](./doc/install_metrics_server.png)

## 启动 etcd-server

- 在 K8S 环境部署好之后, 需要外部启动 etcd 服务

![install etcd](./doc/install_etcd_out_cluster.png)

## 安装 CRD

- 安装 CRD

![install crd](./doc/install_crd.png)

- 安装 RBAC Role
  
![install rbac](./doc/install_rbac.png)

## 本地启动服务

- 启动 disovery 服务

![run-discovery-local](./doc/run-discovery-local.png)

- 启动 authx-apiserver

![run-authx-local](./doc/run-authx-local.png)

- 启动 jingx-apiserver

![run-jingx-local](./doc/run-jingx-local.png)

- 启动 openx-apiserver

![run-openx-local](./doc/run-openx-local.png)

- 启动 authx-frontend 登录认证平台

```sh
make run-authx-frontend-local
```

这个时候登录认证平台之后, 会看到如下界面, 一个是镜像服务, 一个是 K8S 服务

点击后面的箭头进入, 会看到各自的账户、角色和 GVK 管理界面

![authx-dashboard-app](./doc/authx-dashboard-app.png)

进入并复制 openx-apiserver 的登录密码

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
