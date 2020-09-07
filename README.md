# purple

grpc : https://grpc.io/docs/languages/go/quickstart
protoc -I=protobuf --go_out=gen-go --go-grpc_out=gen-go protobuf/member.proto
