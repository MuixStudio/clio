package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/muixstudio/clio/internal/aggregater/handler/auth"
	"github.com/muixstudio/clio/internal/aggregater/handler/health"
	"github.com/muixstudio/clio/internal/aggregater/handler/user"
	"github.com/muixstudio/clio/internal/aggregater/middleware"
	"github.com/muixstudio/clio/internal/aggregater/svc"
)

func Register(r *gin.RouterGroup) {

	svcCtx := svc.NewServiceContext()

	// webcall
	webCallWithoutAuth := r.Group("")
	auth.RegisterWithoutAuth(webCallWithoutAuth, svcCtx)
	health.RegisterWithoutAuth(webCallWithoutAuth, svcCtx)

	webCall := r.Group("")
	webCall.Use(middleware.WebCallAuth())
	auth.Register(webCall, svcCtx)
	user.Register(webCall, svcCtx)

	// api
	api := r.Group("/api")
	user.Register(api, svcCtx)
}
