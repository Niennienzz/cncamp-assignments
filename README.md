# Cloud-Native Camp Assignment #01
<details>
  <summary>点击展开 Assignment #01</summary>

- [GeekBang.org](https://u.geekbang.org/) / [InfoQ.cn](https://www.infoq.cn/) Cloud-Native Camp Assignment #01
- 极客时间云原生训练营 - 作业 #01

## Requirements

- Implement an HTTP server.
- The server handles client requests, and write request headers into response headers.
- The server should read the `VERSION` configuration from the environment, and include it in response headers.
- The server should record client IP & HTTP status code in its log.
- An endpoint `localhost/healthz` should always return 200.

## 要求

- 编写一个 HTTP 服务器。
- 接收客户端 Request，并将 Request 中带的 Header 写入 Response Header -> [middleware/header.go](middleware/header.go)
- 读取当前系统的环境变量中的 `VERSION` 配置，并写入Response Header -> [middleware/header.go](middleware/header.go)
- Server 端记录访问日志包括客户端 IP，HTTP 返回码，输出到 Server 端的标准输出 -> [middleware/logger.go](middleware/logger.go)
- 当访问 `localhost/healthz` 时，应返回 200 OK -> [api/api.go](api/api.go)

## 如何运行

- 测试环境 Golang v1.16+, GNU Make v3.81+.
- 使用 `make dep` 命令下载依赖至 `vendor` 目录。
- 使用 `make test` 命令运行单元测试。
- 可以使用 `make build` 命令编译服务器；也可以使用 `make run` 命令直接运行。
- 由于服务器使用 SQLite，无需创建数据库；运行服务器会默认创建 `sqlite.db` 文件。
- 环境变量与默认值请参考 `config/config.go` 文件。

## 网络请求示例

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

# Cloud-Native Camp Assignment #02
<details>
  <summary>点击展开 Assignment #02</summary>

## 要求

- 构建本地 Docker 镜像
- 编写 `Dockerfile` 将 Assignment #01 编写的服务器容器化 -> [Dockerfile](https://github.com/Niennienzz/cncamp-a01/blob/main/Dockerfile)
- 将镜像推送至 DockerHub 官方镜像仓库
- 通过 Docker 命令本地启动服务器
- 通过 `nsenter` 进入容器查看 IP 配置

## 本地构建与运行

- 构建本地 Docker 镜像

```bash
docker build --tag cncamp_http_server .
```

- 查看镜像列表，成功构建的 `cncamp_http_server` 镜像会出现在列表中

```bash
docker image ls
```

- 通过 Docker 本地启动服务器
- 可以通过 `-e` 传入环境参数

```bash
docker run -p 8080:8080 cncamp_http_server
docker run -p 8080:8080 -e "RATE_LIMIT=5" -e "RATE_LIMIT_WINDOW_SEC=10s" cncamp_http_server
```

## 将镜像推送至 DockerHub

- 镜像已推送至[这里](https://hub.docker.com/repository/docker/niennienzz/cncamp-a02)
- 构建本地 Docker 镜像时打的 Tag 比较简略，推送之前需重新使用标准格式打 Tag

```bash
docker tag <existing-image> <hub-user>/<repo-name>[:<tag>]
```

- 将镜像推送至 DockerHub

```bash
docker push <hub-user>/<repo-name>[:<tag>]
```

## 进入容器查看 IP 配置

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
