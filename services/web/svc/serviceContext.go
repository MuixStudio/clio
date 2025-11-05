package svc

import (
	"context"

	kratosGrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/muixstudio/clio/services/common/pb/userService"
	"github.com/muixstudio/clio/services/web/config"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

type ServiceContext struct {
	Config      config.Config
	RDB         *redis.Client
	UserService userService.UserClient
}

func NewServiceContext(ctx context.Context, config config.Config) (*ServiceContext, error) {
	redisClient, err := newRedisClient(ctx, config)
	grpcClient, err := newGRPCClient(ctx, config)

	userSvc := userService.NewUserClient(grpcClient)

	return &ServiceContext{
		Config:      config,
		RDB:         redisClient,
		UserService: userSvc,
	}, err
}

func newGRPCClient(ctx context.Context, config config.Config) (*grpc.ClientConn, error) {
	grpcClient, err := kratosGrpc.DialInsecure(
		ctx,
		kratosGrpc.WithEndpoint("127.0.0.1:9017"),
	)
	return grpcClient, err
}

func newRedisClient(ctx context.Context, config config.Config) (*redis.Client, error) {
	// init redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Addr,
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
	})

	_, err := rdb.Ping(ctx).Result()
	return rdb, err
}
