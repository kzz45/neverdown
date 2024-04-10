# neverdown

## 本地测试

- 启动 etcd

```sh
docker run -d \
--name etcd-server \
-p 2379:2379 \
-p 2380:2380 \
--env ALLOW_NONE_AUTHENTICATION=yes bitnami/etcd
```

- 启动 discovery

```sh
make run-discovery-local
```

- 启动 authx apiserver

```sh
make run-authx-local
```

- 启动 authx frontend

```sh
make run-authx-frontend-local
```

- 启动 jingx apiserver

```sh
make run-jingx-local
```

- 启动 jingx frontend

```sh
make run-jingx-frontend-local
```

- 启动 openx

```sh
```

## 线上部署
