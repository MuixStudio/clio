package auth

import "github.com/gin-gonic/gin"

type AuthHandler struct {
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (ah AuthHandler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func (ah AuthHandler) Logout() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func (ah AuthHandler) Register() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func (ah AuthHandler) RefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
