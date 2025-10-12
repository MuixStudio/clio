package svc

import (
	"context"
	"fmt"

	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/muixstudio/clio/internal/common/pb/userService"
	"github.com/redis/go-redis/v9"
)

type ServiceContext struct {
	RDB         *redis.Client
	UserService userService.UserClient
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

	userClient, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("127.0.0.1:9017"),
	)

	userSvc := userService.NewUserClient(userClient)

	if err != nil {
		fmt.Printf("failed to new client: %s", err)
	}

	return &ServiceContext{
		RDB:         rdb,
		UserService: userSvc,
	}
}
