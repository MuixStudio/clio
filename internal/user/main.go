package main

import (
	"context"
	"log"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/server"
	"github.com/muixstudio/clio/internal/common/pb/userService/profile"
	"github.com/muixstudio/clio/internal/common/pb/userService/user"
	"github.com/muixstudio/clio/internal/user/config"
	"github.com/muixstudio/clio/internal/user/handler"
	"github.com/muixstudio/clio/internal/user/svc"
)

func main() {

	svr := server.NewServer(
		server.WithErrorHandler(func(ctx context.Context, err error) error {
			klog.CtxErrorf(ctx, "%v", err)
			return err
		}),
		//server.WithTransHandlerFactory(remote.ServerTransHandlerFactory()),
	)

	c := config.Config{}

	svcCtx := svc.NewServiceContext(c)

	err := svr.RegisterService(user.NewServiceInfo(), handler.NewUserImpl(svcCtx))
	err = svr.RegisterService(profile.NewServiceInfo(), new(handler.ProfileImpl))

	klog.SetLevel(klog.LevelDebug)

	err = svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
