package main

import (
	"log"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/muixstudio/clio/internal/common/pb/userService"
	"github.com/muixstudio/clio/internal/user/config"
	"github.com/muixstudio/clio/internal/user/handler"
	"github.com/muixstudio/clio/internal/user/svc"
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
	)
	err := app.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
