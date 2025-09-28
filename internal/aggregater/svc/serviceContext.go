package svc

import (
	"context"
	"fmt"
	"os"

	"github.com/muixstudio/clio/internal/user/pb/user"
	"github.com/redis/go-redis/v9"
	"go-micro.dev/v5"
	"go-micro.dev/v5/transport/grpc"
)

type ServiceContext struct {
	RDB         *redis.Client
	UserService user.UserService
}

func NewServiceContext() *ServiceContext {

	// init redis client
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		//Password: "",
		DB: 0,
	})
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}

	fmt.Println("Redis connect success!")

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
		RDB:         rdb,
		UserService: userService,
	}
}
