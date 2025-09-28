package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/muixstudio/clio/internal/aggregater/svc"
)

func Register(r *gin.RouterGroup, svcCtx *svc.ServiceContext) {

	authHandler := NewAuthHandler(svcCtx)
	r.POST("/logout", authHandler.Logout())

}

func RegisterWithoutAuth(r *gin.RouterGroup, svcCtx *svc.ServiceContext) {

	authHandler := NewAuthHandler(svcCtx)
	r.POST("/refresh_token", authHandler.RefreshToken())
	r.POST("/register", authHandler.Register())
	r.POST("/login", authHandler.Login())
}
