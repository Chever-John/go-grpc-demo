# Pre confirm

预先检验文件。

## 确认 homebrew 的安装路径

在Apple芯片Mac上，Homebrew的安装路径与Intel芯片Mac不同：

```shell
# 检查Homebrew安装位置
which brew

# Apple芯片Mac上应该显示：/opt/homebrew/bin/brew
# Intel芯片Mac上通常是：/usr/local/bin/brew
```

如果尚未安装Homebrew：

```shell
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# 安装完成后，可能需要将Homebrew添加到PATH
echo 'eval "$(/opt/homebrew/bin/brew shellenv)"' >> ~/.zshrc
eval "$(/opt/homebrew/bin/brew shellenv)"
```

## Go 语言环境

```shell
# 检查Go是否已安装及版本
go version

# 如果未安装，通过Homebrew安装（自动支持ARM架构）
brew install go

# 确认安装的是ARM64版本
file $(which go) | grep arm64
```

## Protocol Buffer 编译器（protoc）

```shell
# 检查protoc是否已安装及版本
protoc --version

# 如果未安装，通过Homebrew安装
brew update && brew install protobuf

# 确认安装的是ARM64版本
file $(which protoc) | grep arm64
```

## Go 语言的 Protocol Buffer 插件

```shell
# 安装protoc的Go插件
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# 检查插件是否安装成功
ls $(go env GOPATH)/bin/protoc-gen-go
ls $(go env GOPATH)/bin/protoc-gen-go-grpc
```

## 环境变量设置

```shell
# 检查GOPATH
go env GOPATH

# 确保protoc编译器能找到插件
echo $PATH | grep -q "$(go env GOPATH)/bin" || {
  echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.zshrc
  export PATH="$PATH:$(go env GOPATH)/bin"
}

# 确认PATH设置
echo $PATH
```

## 验证 gRPC 依赖

```shell
# 创建临时目录验证依赖
mkdir -p /tmp/grpc-test && cd /tmp/grpc-test

# 初始化Go模块
go mod init grpctest

# 安装gRPC依赖
go get -u google.golang.org/grpc
go get -u google.golang.org/protobuf/proto

# 检查依赖是否安装成功
go list -m google.golang.org/grpc
```

## 确认权限设置

Apple芯片Mac上有时会遇到权限问题：

```shell
# 确认GOPATH目录权限
ls -la $(go env GOPATH)

# 如果有权限问题，设置正确权限
sudo chown -R $(whoami):staff $(go env GOPATH)
```

## 测试完整流程

```shell
# 创建测试目录
mkdir -p /tmp/grpc-test/proto
cd /tmp/grpc-test

# 创建测试proto文件
cat > proto/test.proto << EOF
syntax = "proto3";
option go_package = "grpctest/proto";
package proto;

service TestService {
  rpc SayHello (HelloRequest) returns (HelloResponse) {}
}

message HelloRequest {
  string name = 1;
}

message HelloResponse {
  string message = 1;
}
EOF

# 生成Go代码
protoc --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  proto/test.proto

# 检查生成的文件
ls proto/test*.go
```

## 验证gRPC服务器和客户端

创建简单的服务器和客户端代码测试gRPC连接：

```shell
# 创建服务器代码
mkdir -p /tmp/grpc-test/server
cat > /tmp/grpc-test/server/main.go << EOF
package main

import (
	"context"
	"log"
	"net"
	
	pb "grpctest/proto"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedTestServiceServer
}

func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Printf("收到请求: %v", req.GetName())
	return &pb.HelloResponse{Message: "你好，" + req.GetName()}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("监听失败: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterTestServiceServer(s, &server{})
	log.Printf("服务器在 %v 上运行", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("服务失败: %v", err)
	}
}
EOF

# 创建客户端代码
mkdir -p /tmp/grpc-test/client
cat > /tmp/grpc-test/client/main.go << EOF
package main

import (
	"context"
	"log"
	"time"
	
	pb "grpctest/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	defer conn.Close()
	
	c := pb.NewTestServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: "云原生开发者"})
	if err != nil {
		log.Fatalf("调用失败: %v", err)
	}
	log.Printf("响应: %s", r.GetMessage())
}
EOF
```

如果完成上述所有步骤且没有错误，你的Apple芯片Mac上的Golang和gRPC环境已经配置成功，可以开始开发gRPC应用了 [1](https://grpc.io/docs/languages/go/quickstart/)[2](https://abc101.medium.com/golang-grpc-mac-2fe01939a29d).