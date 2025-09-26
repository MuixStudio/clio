package testclient

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/muixstudio/clio/internal/user/pb/user"
	"go-micro.dev/v5"
	"go-micro.dev/v5/transport/grpc"
)

func TestClient(t *testing.T) {
	// 在初始化 service 之前临时清理 go test 注入的标志，避免 go-micro 解析失败
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

	client := user.NewUserService("user.User", service.Client())
	//phone := ""
	//name := ""
	username := ""
	password := ""
	//isAdmin := true
	rsp, err := client.CreateUser(context.Background(), &user.CreateUserRequest{
		//Name:     &name,
		Password: &password,
		UserName: &username,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(rsp)
}
