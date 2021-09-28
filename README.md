# Cloud-Native Camp Assignment #01

- [GeekBang.org](https://u.geekbang.org/) / [InfoQ.cn](https://www.infoq.cn/) Cloud-Native Camp Assignment #01
- 极客时间云原生训练营 - 作业 #01

## Requirements

- Implement an HTTP server.
- The server handles client requests, and write request headers into response headers.
- The server should read the `VERSION` configuration from the environment, and include it in response headers.
- The server should record client IP, HTTP status code in its log.
- An endpoint `localhost/healthz` should always return 200.

## 要求

- 编写一个 HTTP 服务器
- 接收客户端 Request，并将 Request 中带的 Header 写入 Response Header
- 读取当前系统的环境变量中的 `VERSION` 配置，并写入Response Header
- Server 端记录访问日志包括客户端 IP，HTTP 返回码，输出到 Server 端的标准输出
- 当访问 localhost/healthz 时，应返回 200
