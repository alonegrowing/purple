# Purple 基于 gin & grpc 的 web/rpc 一站式服务框架

- 支持配置的管理、解析
- 支持 redis / mysql 的连接池
- 支持 gorm 的便捷 mysql 数据库操作
- 支持 资源操作的打点监控

- grpc : https://grpc.io/docs/languages/go/quickstart
- protoc -I=proto --go_out=gen-go --go-grpc_out=gen-go proto/*.proto

# required
- go version go1.14.8 darwin/amd64
- libprotoc 3.13.0


