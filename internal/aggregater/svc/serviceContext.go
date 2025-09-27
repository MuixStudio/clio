package svc

import (
	"os"

	"github.com/muixstudio/clio/internal/user/pb/user"
	"go-micro.dev/v5"
	"go-micro.dev/v5/transport/grpc"
)

type ServiceContext struct {
	UserService user.UserService
}

func NewServiceContext() *ServiceContext {
	grpcTransport := grpc.NewTransport()
	origArgs := os.Args
	defer func() { os.Args = origArgs }()
	os.Args = []string{"user.client"}

	service := micro.NewService(
		micro.Name("user.client"),
		micro.Version("0.0.1"),
	)
	service.Init(
		micro.Transport(grpcTransport),
	)

	userService := user.NewUserService("user.User", service.Client())
	return &ServiceContext{
		UserService: userService,
	}
}
