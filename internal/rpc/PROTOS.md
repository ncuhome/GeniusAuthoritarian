**pull git submodules and enter this dir first**

#### app.proto

```shell
protoc --go_out=.\app\appProto --go-grpc_out=.\app\appProto .\protos\app.proto
```

#### refreshToken.proto

```shell
protoc --go_out=.\refreshToken\refreshTokenProto --go-grpc_out=.\refreshToken\refreshTokenProto .\protos\refreshToken.proto
```