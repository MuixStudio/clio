package main

import (
	"context"
	"fmt"
	"log"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/muixstudio/clio/services/common/pb/userService"
	"github.com/muixstudio/clio/services/user/config"
	"github.com/muixstudio/clio/services/user/handler"
	"github.com/muixstudio/clio/services/user/svc"
)

func main() {

	c := config.Config{}

	svcCtx := svc.NewServiceContext(c)

	gs := grpc.NewServer(
		grpc.Address(":9017"),
	)

	userService.RegisterUserServer(gs, handler.NewUserHandler(svcCtx))
	app := kratos.New(
		kratos.Name("user"),
		kratos.Version("v1.0.0"),
		kratos.Server(gs),
		kratos.AfterStart(func(ctx context.Context) error {
			fmt.Println("user service started")
			return nil
		}),
	)
	err := app.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
