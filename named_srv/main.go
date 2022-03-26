package main

import (
	"context"
	"fmt"
	"net"
	"strings"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"

	"github.com/lpxxn/grpc_demo/common"
	"github.com/lpxxn/grpc_demo/protos/api"
	"github.com/lpxxn/grpc_demo/protos/model"
)

func main() {
	ip, err := common.PrivateIPv4()
	if err != nil {
		panic(err)
	}

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}
	addr := fmt.Sprintf("%s:%d", ip.String(), listener.Addr().(*net.TCPAddr).Port)
	fmt.Println(addr)

	s := grpc.NewServer([]grpc.ServerOption{
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
				fmt.Printf("current service addr: %s \n", addr)
				fmt.Println(strings.Repeat("-", 20))
				return handler(ctx, req)
			},
		)),
	}...)
	api.RegisterStudentSrvServer(s, &common.StudentSrv{StudentList: []*model.Student{
		&model.Student{Id: 1, Name: "tom", Age: 5},
		&model.Student{Id: 2, Name: "jerry", Age: 6},
	}})

	reg, err := common.NewService(common.ServiceInfo{
		Name: common.ServName,
		Addr: addr,
	}, []string{"http://127.0.0.1:2379"})
	if err != nil {
		panic(err)
	}
	go func() {
		if err := reg.Start(); err != nil {
			panic(err)
		}
	}()

	if err := s.Serve(listener); err != nil {
		fmt.Println(err)
	}
}
