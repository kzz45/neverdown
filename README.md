# neverdown

## 说明

一个 K8S 管理平台, 包含权限管理, 镜像管理, K8S 集群管理

已经验证阿里云、腾讯云、火山云上的 K8S 服务

## 模块说明

| 模块            | 说明                 |
| :-------------- | :------------------- |
| authx-frontend  | 权限管理, 前端       |
| jingx-frontend  | 镜像管理, 前端       |
| openx-apiserver | K8S 集群的管理       |
| openx-frontend  | K8S 集群的管理, 前端 |

## 安装部署

```sh
kubectl apply -f config/crd/

kubectl apply -f config/deploy/
```

有任何问题, 请联系我:

![wechat](./doc/Wechat.png)
