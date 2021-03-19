# halfway

[![Language](https://img.shields.io/badge/Language-Go-blue.svg)](https://golang.org/)
![Go Version](https://camo.githubusercontent.com/7b5b1d36152e872a911a5e21c1aad097f74c0292960edc143923a3b2fbe5f458/68747470733a2f2f696d672e736869656c64732e696f2f62616467652f676f25323076657273696f6e2d2533453d312e31342d3631434644442e7376673f7374796c653d666c61742d737175617265)

整合golang开源框架搭建的微服务，旨在快速搭建稳定高效的开发平台

## Features
- HTTP Server：基于[echo](https://github.com/labstack/echo)框架设计，集成 
  限流 熔断 请求日志 中文验证器 等中间件， 参考[jupiter](https://github.com/douyu/jupiter).GRPCProxyWrapper封装grpc.UnaryHandler转化为echo.HandleFunc；
- RPC Server：基于官方gRPC开发，集成discovery模块支持服务发现功能，基于[gogo/protobuf](https://github.com/gogo/protobuf)生成pb.go文件；
- Cache：Redis封装[go-redis](https://github.com/go-redis/redis/v8)框架；
- DB：Mysql采用[gorm](https://gorm.io)框架，ES，Tidb，Hbase还待开发；
- Config：采用[viper](https://github.com/spf13/viper)，配合consul K/V 作为远程配置中心，可实现配置版本管理和更新；
- Log：采用[zap](https://github.com/uber-go/zap)的field实现高性能日志库，并结合 filebeat elk 实现远程日志管理；
- 全链路trace基于[elastic APM](https://www.elastic.co/guide/en/apm/agent/go/current/index.html)，支持(gRPC/HTTP/MySQL/Redis)，集群接入linkerd服务网格化后逐步替换；

## Demo Server
https://github.com/xulichen/halfway_demo


## File Structure
```text
├── README.md
├── doc
│   ├── example.yaml            
│   └── version.md
├── go.mod
├── go.sum
├── pkg
│   ├── cache
│   │   └── redis.go
│   ├── config
│   │   ├── base.go
│   │   ├── config.go
│   │   ├── config_test.go
│   │   └── redis.go
│   ├── consts
│   │   ├── consts.go
│   │   └── errors.go
│   ├── db
│   │   └── mysql.go
│   ├── discovery
│   │   ├── consul
│   │   │   ├── consul.go
│   │   │   └── resolver.go
│   │   └── discovery.go
│   ├── log
│   │   └── log.go
│   ├── net
│   │   ├── http
│   │   │   ├── client.go
│   │   │   ├── config.go
│   │   │   ├── errors
│   │   │   │   └── error.go
│   │   │   ├── http.go
│   │   │   ├── middleware
│   │   │   │   ├── circuit_breaker.go
│   │   │   │   ├── echo_logger.go
│   │   │   │   ├── rate_limit.go
│   │   │   │   ├── recover.go
│   │   │   │   └── validate.go
│   │   │   └── server.go
│   │   └── rpc
│   │       ├── client.go
│   │       ├── config.go
│   │       ├── middleware
│   │       │   └── validator.go
│   │       └── server.go
│   └── utils
│       ├── common.go
│       └── validator.go
├── test
│   └── test.md
├── third_party
└── tools
    └── tools.md

```

## TODO
- 基于cobra开发tools工具包，支持命令行生成demo文件
- 基于protobuf文件生成demo
- 项目优化，完成@@TODO