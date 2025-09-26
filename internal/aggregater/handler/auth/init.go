package auth

import "github.com/gin-gonic/gin"

func Register(r *gin.RouterGroup) {
	r.POST("/login", NewAuthHandler().Login())
	r.POST("/logout", NewAuthHandler().Logout())
	r.POST("/register", NewAuthHandler().Register())
	r.POST("/refresh_token", NewAuthHandler().RefreshToken())
}
