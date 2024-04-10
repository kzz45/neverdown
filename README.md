# neverdown

## 本地测试

启动 etcd

```sh
docker run -d \
--name etcd-server \
-p 2379:2379 \
-p 2380:2380 \
--env ALLOW_NONE_AUTHENTICATION=yes bitnami/etcd
```

启动 discovery

启动 authx

启动 jingx

启动 openx

## 线上部署
