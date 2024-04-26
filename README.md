# neverdown

## 说明

一个 K8S 管理平台, 包含权限管理, 镜像管理, K8S 集群管理

## 模块说明

| 模块            | 说明                 |
| :-------------- | :------------------- |
| discovery       | 服务发现             |
| authx-apiserver | 权限管理             |
| authx-frontend  | 权限管理, 前端       |
| jingx-apiserver | 镜像管理             |
| jingx-frontend  | 镜像管理, 前端       |
| openx-apiserver | K8S 集群的管理       |
| openx-frontend  | K8S 集群的管理, 前端 |

## [本地测试](./local.md)

## 安装部署

```sh
kubectl apply -f config/crd/

kubectl apply -f config/deploy/
```

## 联系

有任何问题, 请联系我:

![wechat](./doc/Wechat.jpg)