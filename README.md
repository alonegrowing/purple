# purple

grpc : https://grpc.io/docs/languages/go/quickstart
protoc -I=protobuf --go_out=gen-go --go-grpc_out=gen-go protobuf/*.proto






protoc
  --go_out=Mgrpc/service_config/service_config.proto=/internal/proto/grpc_service_config:. \
  --go-grpc_out=Mgrpc/service_config/service_config.proto=/internal/proto/grpc_service_config:. \
  --go_opt=paths=source_relative \
  --go-grpc_opt=paths=source_relative \
  helloworld/helloworld.proto