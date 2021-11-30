# Cloud-Native Camp Assignment #01 - Go
<details>
  <summary>Click to expand Assignment #01</summary>

- [GeekBang.org](https://u.geekbang.org/) / [InfoQ.cn](https://www.infoq.cn/) Cloud-Native Camp Assignment #01
- 极客时间云原生训练营 - 作业 #01

## 1.1 - Requirements

- Implement an HTTP server.
- The server handles client requests, and write request headers into response headers.
- The server should read the `VERSION` configuration from the environment, and include it in response headers.
- The server should record client IP & HTTP status code in its log.
- An endpoint `localhost/healthz` should always return 200.

## 1.2 - 要求

- 编写一个 HTTP 服务器。
- 接收客户端 Request，并将 Request 中带的 Header 写入 Response Header -> [middleware/header.go](middleware/header.go)
- 读取当前系统的环境变量中的 `VERSION` 配置，并写入Response Header -> [middleware/header.go](middleware/header.go)
- Server 端记录访问日志包括客户端 IP，HTTP 返回码，输出到 Server 端的标准输出 -> [middleware/logger.go](middleware/logger.go)
- 当访问 `localhost/healthz` 时，应返回 200 OK -> [api/api.go](api/api.go)

## 1.3 - 如何运行

- 测试环境 Golang v1.16+, GNU Make v3.81+.
- 使用 `make dep` 命令下载依赖至 `vendor` 目录。
- 使用 `make test` 命令运行单元测试。
- 可以使用 `make build` 命令编译服务器；也可以使用 `make run` 命令直接运行。
- 由于服务器使用 SQLite，无需创建数据库；运行服务器会默认创建 `sqlite.db` 文件。
- 环境变量与默认值请参考 `config/config.go` 文件。

## 1.4 - 网络请求示例

- Healthz 检测

  ```bash
  curl --request GET --url http://localhost:8080/healthz
  ```

- 用户注册
- 邮箱地址格式须合法，密码长度至少八位

  ```bash
  curl --request POST --url http://localhost:8080/user/signup \
       --header 'Content-Type: application/json' \
       --data '{
         "email": "someuser@test.com",
         "password": "12345678"
       }'
  ```

- 用户登录
- 成功登录后获取 `{accessToken}`

  ```bash
  curl --request POST --url http://localhost:8080/user/login \
       --header 'Content-Type: application/json' \
       --data '{
         "email": "someuser@test.com",
         "password": "12345678"
       }'
  ```

- Crypto 价格
- 因为是简单示例服务器，路径参数 `{cryptoCode}` 仅支持 `ADA`、`BNB`、`BTC` 与 `ETH`
- Authorization Header 需要使用上述登录获取的 `{accessToken}`

  ```bash
  curl --request GET --url http://localhost:8080/crypto/{cryptoCode} \
       --header 'Authorization: Bearer {accessToken}'
  ```

</details>

---

# Cloud-Native Camp Assignment #02 - Docker
<details>
  <summary>Click to expand Assignment #02</summary>

## 2.1 - 要求

- 构建本地 Docker 镜像
- 编写 `Dockerfile` 将 Assignment #01 编写的服务器容器化 -> [Dockerfile](https://github.com/Niennienzz/cncamp-a01/blob/main/Dockerfile)
- 将镜像推送至 DockerHub 官方镜像仓库
- 通过 Docker 命令本地启动服务器
- 通过 `nsenter` 进入容器查看 IP 配置

## 2.2 - 本地构建与运行

- 构建本地 Docker 镜像

  ```bash
  make image
  ```

- 查看镜像列表，成功构建的 `niennienzz/cncamp_http_server` 镜像会出现在列表中

  ```bash
  docker image ls
  ```

- 通过 Docker 本地启动服务器
- 可以通过 `-e` 传入环境参数

  ```bash
  docker run -p 8080:8080 cncamp_http_server
  docker run -p 8080:8080 -e "RATE_LIMIT=5" -e "RATE_LIMIT_WINDOW_SEC=10s" cncamp_http_server
  ```

## 2.3 - 将镜像推送至 DockerHub

- 镜像已推送至[这里](https://hub.docker.com/repository/docker/niennienzz/cncamp_http_server)
- 将镜像推送至 DockerHub

  ```bash
  make push
  ```

## 2.3 - 进入容器查看 IP 配置

- 找到正在运行的容器实例

  ```bash
  docker ps | grep cncamp_http_server
  #=> 6592fd79xxxx
  ```

- 找到上述容器实例实例的 PID

  ```bash
  docker inspect 6592fd79xxxx | grep -i pid
  #=> 12345
  ```

- 通过 `nsenter` 进入容器查看 IP 配置

  ```bash
  nsenter -t 12345 -n ip a
  ```

</details>

---

# Cloud-Native Camp Assignment #03 - Kubernetes

## 3.1 - 文件路径结构

- 文件路径:
  > 代码均在 `/httpserver` 路径中。
  >
  > 配置均在 `/deployment` 路径中。

- 在 `/deployment` 路径中，对于某个服务 `service` 来说:
  > `Service` & `Deployment` 均集中配置在 `{service}.yaml` 文件中。
  >
  > `ConfigMap` 均配置在 `{service}-config.yaml` 文件中。
  >
  > `Secret` 均配置在 `{service}-secret.yaml` 文件中。
  >
  > 遵循命名约定，以此类推。

## 3.2 - 要求与分析

### 编写 Kubernetes 部署脚本将 `httpserver` 部署到 Kubernetes 集群

- 优雅启动

  > 使用 `readinessProbe` 探针检查 Pod 是否就绪，就绪则接收请求。
  >
  > 查看 `httpserver.yaml` 文件中的 `readinessProbe` 部分。

- 优雅终止

  > 使用 `terminationGracePeriodSeconds` 给与 Pod 适当的终止时间。
  >
  > 查看 `/deployment/httpserver.yaml` 文件中的 `terminationGracePeriodSeconds` 部分。
  >
  > 当 Pod 关闭 Kubernetes 将给应用发送 `SIGTERM` 信号并等待配置的时间后关闭。
  >
  > 同时 `httpserver` 代码中接受 `SIGTERM` 信号并执行各项终止任务，例如关闭数据库连接等。

- 资源需求和服务质量保证

  > 查看 `/deployment/httpserver.yaml` 文件中的 `resources` 部分。

- 探活

  > 使用 `livenessProbe` 探针检查 Pod 是否存活，如果检测不到 Pod 存活则杀掉当前 Pod 重启。
  >
  > 查看 `/deployment/httpserver.yaml` 文件中的 `livenessProbe` 部分。

- 日常运维需求，日志等级

  > `httpserver` 代码中使用 [`logrus`](https://github.com/sirupsen/logrus) 库中不同的日志级别。

- 配置和代码分离

  > 代码全部在 `/httpserver` 路径中，而配置全部在 `/deployment` 路径中。
  >
  > 利用配置中的 `*-config.yaml` & `*-secret.yaml` 文件将配置注入到 Pod 中。

- 如何确保整个应用的高可用

  > 首先保证 `httpserver` 代码本身是完全无状态的，例如没有本地内存缓存；状态均集中保存在数据库中。
  >
  > 配置增加多个副本，查看 `/deployment/httpserver.yaml` 文件中的 `replicas` 部分。

- 如何通过证书保证通讯安全

  > `httpserver` 代码本身没有使用 HTTPS 证书。
  >
  > 在配置 Ingress 时使用证书保证通讯安全。

## 3.3 - 实验

### 3.3.1 - 实验环境

- 因为没有申请服务器，且本机配置还不错，所以采用本地 Minikube 进行实验。
- 本地实验肯定是规避了一些远程操作的困难，希望自己在以后的练习当中多尝试，而非浅尝辄止。

### 3.3.2 - 实验准备

- 生成 HTTPS 证书

  ```bash
  openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout tls.key -out tls.crt -subj "/CN=niennienzz-httpserver.com/O=niennienzz-httpserver"
  ```

- 生成 `/deployment/httpserver-tls-secret.yaml` 文件

  ```bash
  kubectl create secret tls tls-secret --cert=tls.crt --key=tls.key --dry-run=client -o yaml
  ```

- 预先安装 Minikube Ingress Controller 组件

  ```bash
  minikube addons enable ingress
  ```

### 3.3.3 - 实验步骤

使用 `make cluster` 创建集群；使用 `make destroy` 销毁集群；下面是详细步骤解析。

- 配置 ConfigMap 与 Secret

  ```bash
  kubectl apply -f deployment/pv.yaml
  kubectl apply -f deployment/pvc.yaml
  kubectl apply -f deployment/mongo-config.yaml
  kubectl apply -f deployment/mongo-secret.yaml
  kubectl apply -f deployment/httpserver-config.yaml
  kubectl apply -f deployment/httpserver-secret.yaml
  kubectl apply -f deployment/httpserver-tls-secret.yaml
  ```

- 配置 Service

  ```bash
  kubectl apply -f deployment/mongo.yaml
  kubectl apply -f deployment/httpserver.yaml
  ```

- 配置 Ingress Rules

  ```bash
  kubectl apply -f deployment/httpserver-ingress.yaml
  ```

### 3.3.4 - 实验结果

- 查看 Pod、Service、Deployment

  ```bash
  kubectl get all
  #=> NAME                                         READY   STATUS    RESTARTS      AGE
  #=> pod/httpserver-deployment-7dc79b84f8-4jzpw   1/1     Running   0             58m
  #=> pod/httpserver-deployment-7dc79b84f8-g82dx   1/1     Running   0             58m
  #=> pod/httpserver-deployment-7dc79b84f8-tstxb   1/1     Running   0             58m
  #=> pod/mongo-deployment-7875498c-8tbch          1/1     Running   9 (27h ago)   28h
  
  #=> NAME                         TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)          AGE
  #=> service/httpserver-service   NodePort    10.108.184.166   <none>        8080:30100/TCP   58m
  #=> service/kubernetes           ClusterIP   10.96.0.1        <none>        443/TCP          45h
  #=> service/mongo-service        ClusterIP   10.102.201.39    <none>        27017/TCP        28h
  
  #=> NAME                                    READY   UP-TO-DATE   AVAILABLE   AGE
  #=> deployment.apps/httpserver-deployment   3/3     3            3           58m
  #=> deployment.apps/mongo-deployment        1/1     1            1           28h
  
  #=> NAME                                               DESIRED   CURRENT   READY   AGE
  #=> replicaset.apps/httpserver-deployment-7dc79b84f8   3         3         3       58m
  #=> replicaset.apps/mongo-deployment-7875498c          1         1         1       28h
  ```

- 查看 Ingress

  ```bash
  kubectl get ingress
  #=> NAME                 CLASS   HOSTS                       ADDRESS     PORTS     AGE
  #=> httpserver-ingress   nginx   niennienzz-httpserver.com   localhost   80, 443   68m
  ```

  ```bash
  kubectl get svc -n ingress-nginx
  #=> NAME                                 TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)                      AGE
  #=> ingress-nginx-controller             NodePort    10.98.165.142   <none>        80:30090/TCP,443:31729/TCP   89m
  #=> ingress-nginx-controller-admission   ClusterIP   10.111.53.73    <none>        443/TCP                      89m
  ```

- 查看 Minikube 地址

  ```bash
  minikube ip
  #=> 192.xxx.xx.x
  ```

- 修改 `/etc/hosts` 文件添加 `{minikube_ip} niennienzz-httpserver.com` 之后集群可以接收请求

  ```bash
  curl --insecure --request POST \
    --url https://niennienzz-httpserver.com/user/signup \
    --header 'Content-Type: application/json' \
    --data '{
      "email": "someuser_01@test.com",
      "password": "12345678"
  }'
  ```

  ```bash
  curl --insecure --request POST \
    --url https://niennienzz-httpserver.com/user/login \
    --header 'Content-Type: application/json' \
    --data '{
      "email": "someuser_01@test.com",
      "password": "12345678"
  }'
  ```

  ```bash
  curl --insecure --request GET \
    --url http://192.168.64.2:30100/crypto/ETH \
    --header 'Authorization: Bearer {TOKEN}'
  ```