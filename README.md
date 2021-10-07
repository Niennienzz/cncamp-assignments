# Cloud-Native Camp Assignment #01

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
