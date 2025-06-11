package main

import (
	"context"
	"flag"
	"fmt"
	demo "github.com/Chever-John/go-grpc-demo/pkg/pb/v1"
	"io"
	"log"
	"net"
	"time"

	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"

	"google.golang.org/grpc"
)

type demoServer struct {
	demo.UnimplementedDemoServer
	savedResults []*demo.Response //用于服务端流
}

// contextWithTimeout 封装 Context 超时处理
func contextWithTimeout(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	if _, ok := ctx.Deadline(); !ok {
		return context.WithTimeout(ctx, timeout)
	}
	return ctx, func() {} // 如果已经有 Deadline，则不创建新的，返回一个空的 CancelFunc
}

// Add 实现方法Add
func (s *demoServer) Add(ctx context.Context, in *demo.TwoNum) (*demo.Response, error) {
	// use the deadline from ctx
	ctx, cancel := contextWithTimeout(ctx, 5*time.Second)
	defer cancel()

	x := in.X
	y := in.Y

	result := x + y

	return &demo.Response{
		Result: result,
	}, nil
}

// SayHello 实现方法SayHello
func (s *demoServer) SayHello(ctx context.Context, in *demo.HelloRequest) (*demo.HelloReply, error) {
	ctx, cancel := contextWithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 1. check nil pointer
	if in == nil {
		return nil, fmt.Errorf("HelloRequest is nil")
	}
	// 2. check Name field, defense it is nil
	name := in.GetName()
	if name == "" {
		return nil, fmt.Errorf("name is empty")
	}

	// 3. use the safer string
	message := "Hello " + name

	return &demo.HelloReply{Message: message}, nil
}

// GetStream 实现方法GetStream
func (s *demoServer) GetStream(in *demo.TwoNum, pipe demo.Demo_GetStreamServer) error {

	err := pipe.Send(&demo.Response{Result: in.X + in.Y}) //返回和
	if err != nil {
		return err
	}
	err = pipe.Send(&demo.Response{Result: in.X * in.Y}) //返回积
	if err != nil {
		return err
	}
	err = pipe.Send(&demo.Response{Result: in.X - in.Y}) //返回差
	if err != nil {
		return err
	}

	return nil
}

// PutStream 实现方法PutStream
func (s *demoServer) PutStream(pipe demo.Demo_PutStreamServer) error {
	var res int32
	for { //循环接收
		request, err := pipe.Recv()
		if err == io.EOF { //判断是否发送结束
			break
		}
		if err != nil {
			log.Println(err.Error())
		}

		if request == nil {
			log.Println("Received request is nil")
			continue
		}
		res += request.X //累加
	}
	_ = pipe.SendAndClose(&demo.Response{Result: res}) //返回
	return nil
}

// DoubleStream 实现方法DoubleStream
func (s *demoServer) DoubleStream(pipe demo.Demo_DoubleStreamServer) error {

	for {
		request, err := pipe.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if err = pipe.Send(&demo.Response{Result: request.X + request.Y}); err != nil {
			return err
		}

	}

}

// SendLargeData 处理大型数据请求，可能会触发 frame too large 错误
func (s *demoServer) SendLargeData(ctx context.Context, req *demo.LargeRequest) (*demo.LargeResponse, error) {
	payloadSize := len(req.LargePayload)
	log.Printf("Received large payload request, and the length of it is %d", payloadSize)

	return &demo.LargeResponse{
		Status:      "Have received data",
		PayloadSize: int32(payloadSize),
	}, nil
}

var (
	tls        = flag.Bool("tls", false, "使用启用tls") //默认false
	port       = flag.Int("port", 50054, "服务端口")    //默认50054
	maxMsgSize = flag.Int("max_msg_size", 4*1024*1024, "Max message size(bytes)")
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalln(err)
	}

	var opts []grpc.ServerOption

	if *tls {
		creds, err := credentials.NewServerTLSFromFile("certs/server/server.crt", "certs/server/server.key")
		if err != nil {
			log.Fatalf("Failed to generate credentials %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}

	// 默认不设置较大的消息大小，方便复现错误
	// 如果要修复错误，可以取消下面两行的注释
	// opts = append(opts, grpc.MaxRecvMsgSize(*maxMsgSize))
	// opts = append(opts, grpc.MaxSendMsgSize(*maxMsgSize))

	s := grpc.NewServer(opts...)
	demo.RegisterDemoServer(s, &demoServer{})
	reflection.Register(s)
	log.Printf("Server listeing at :%v, Max message: %v bytes\n", *port, *maxMsgSize)
	if err := s.Serve(lis); err != nil {
		log.Fatalln(err)
	}
}
