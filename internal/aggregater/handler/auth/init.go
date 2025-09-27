package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/muixstudio/clio/internal/aggregater/svc"
)

func Register(r *gin.RouterGroup, svcCtx *svc.ServiceContext) {

	authHandler := NewAuthHandler(svcCtx)

	r.POST("/login", authHandler.Login())
	r.POST("/logout", authHandler.Logout())
	r.POST("/register", authHandler.Register())
	r.POST("/refresh_token", authHandler.RefreshToken())
}
