package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/muixstudio/clio/internal/aggregater/handler/auth"
	"github.com/muixstudio/clio/internal/aggregater/handler/user"
	"github.com/muixstudio/clio/internal/aggregater/svc"
)

func Register(r *gin.RouterGroup) {

	svcCtx := svc.NewServiceContext()
	auth.Register(r, svcCtx)
	user.Register(r, svcCtx)
}
