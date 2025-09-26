package main

import (
	"github.com/gin-gonic/gin"
	"github.com/muixstudio/clio/internal/aggregater/handler"
)

//type LoginRequest struct {
//	Username    *string `json:"username,omitempty" binding:"required"`
//	Name        *string `json:"name,omitempty" binding:"required"`
//	Password    *string `json:"password,omitempty" binding:"required"`
//	IsAdmin *bool   `json:"is_admin,omitempty"`
//	CountryCode *string `json:"country_code,omitempty"`
//	Phone *string `json:"phone,omitempty"`
//	Email *string `json:"email,omitempty"`
//}

func main() {

	//grpcTransport := grpc.NewTransport()
	//origArgs := os.Args
	//defer func() { os.Args = origArgs }()
	//os.Args = []string{"user.client"}
	//
	//service := micro.NewService(
	//	micro.Name("user.client"),
	//	micro.Version("0.0.1"),
	//)
	//service.Init(
	//	micro.Transport(grpcTransport),
	//)
	//
	//userService := user.NewUserService("user.User", service.Client())

	r := gin.Default()
	handler.Register(&r.RouterGroup)
	if err := r.Run(":5020"); err != nil {
		panic(err)
	}
}
