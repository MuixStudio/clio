package svc

import (
	"context"
	"fmt"

	"github.com/cloudwego/kitex/client"
	"github.com/muixstudio/clio/internal/common/pb/userService/user"
	"github.com/redis/go-redis/v9"
)

type ServiceContext struct {
	RDB         *redis.Client
	UserService user.Client
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

	//grpcTransport := grpc.NewTransport()
	//origArgs := os.Args
	//defer func() { os.Args = origArgs }()
	//os.Args = []string{"user.client"}
	//
	//service := micro.NewService(
	//	micro.Name("aggregater.client"),
	//	micro.Version("0.0.1"),
	//)
	//service.Init(
	//	micro.Transport(grpcTransport),
	//)
	//
	//c := service.Client()
	//
	//err = c.Init(
	//	client.Transport(grpcTransport),
	//	client.PoolTTL(time.Second*20),
	//	client.PoolSize(11),
	//	client.PoolCloseTimeout(time.Second*10),
	//	client.DialTimeout(time.Second*10),
	//)
	//if err != nil {
	//	panic(err)
	//}

	userClient, err := user.NewClient("userServiceInfo", client.WithHostPorts("127.0.0.1:8888"))

	if err != nil {
		fmt.Printf("failed to new client: %s", err)
	}

	userService := userClient
	return &ServiceContext{
		RDB:         rdb,
		UserService: userService,
	}
}
