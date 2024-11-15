# gRPC 
```
go get -u google.golang.org/grpc
go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
```

protoc コマンド

```
protoc --go_out=pkg/grpc/storage --go_opt=paths=source_relative --go-grpc_out=pkg/grpc/storage --go-grpc_opt=paths=source_relative proto/storage.proto
```