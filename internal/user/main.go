package main

import (
	"log"

	"github.com/muixstudio/clio/internal/user/config"
	"github.com/muixstudio/clio/internal/user/handler"
	"github.com/muixstudio/clio/internal/user/pb/user"
	"go-micro.dev/v5"
	"go-micro.dev/v5/transport/grpc"
)

func main() {
	grpcTransport := grpc.NewTransport()

	// 创建服务
	service := micro.NewService(
		micro.Name("user.User"),
		micro.Version("0.0.1"),
	)

	// 初始化
	service.Init(
		micro.Transport(grpcTransport),
	)

	// 注册 handler
	var c config.Config
	err := user.RegisterUserHandler(service.Server(), handler.NewUserHandler(c))

	if err != nil {
		log.Fatal(err)
	}

	// 运行服务
	if err = service.Run(); err != nil {
		log.Fatal(err)
	}
}
