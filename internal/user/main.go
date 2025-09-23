package main

import (
	"log"

	"github.com/muixstudio/clio/internal/user/handler"
	"github.com/muixstudio/clio/internal/user/pb/user"
	"go-micro.dev/v5"
	"go-micro.dev/v5/transport/grpc"
)

func main() {
	grpcTransport := grpc.NewTransport()

	// 创建服务
	service := micro.NewService(
		micro.Name("user.prc"),
		micro.Version("0.0.1"),
	)

	// 初始化
	service.Init(
		micro.Transport(grpcTransport),
	)

	// 注册 handler
	user.RegisterUserHandler(service.Server(), handler.NewUserHandler())

	// 运行服务
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
