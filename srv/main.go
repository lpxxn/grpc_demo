package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"

	"github.com/lpxxn/grpc_demo/common"
	_ "github.com/lpxxn/grpc_demo/mycodec"
	"github.com/lpxxn/grpc_demo/protos/api"
	"github.com/lpxxn/grpc_demo/protos/model"
	"google.golang.org/grpc/reflection"
)

var port int

func init() {
	flag.IntVar(&port, "port", 10001, "grpc port")
}

func main() {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(WithInterceptor()...)
	api.RegisterStudentSrvServer(grpcServer, &common.StudentSrv{StudentList: []*model.Student{
		&model.Student{Id: 1, Value: "tom", Age: 5},
		&model.Student{Id: 2, Value: "jerry", Age: 6},
		//&model.Student{Id: 1, Value: "tom", Age: 5},
		//&model.Student{Id: 2, Value: "jerry", Age: 6},
	}, Version: "simple srv"})
	// 生产环境就不要注册了
	reflection.Register(grpcServer)
	fmt.Println("serv running...")
	if err := grpcServer.Serve(listen); err != nil {
		panic(err)
	}
}

func WithInterceptor() []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
				// log metric ...
				fmt.Printf("%+v\n", info)
				return handler(ctx, req)
			},
		)),
	}
}
