package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/muixstudio/clio/internal/aggregater/handler/auth"
	"github.com/muixstudio/clio/internal/aggregater/handler/user"
)

func Register(r *gin.RouterGroup) {
	auth.Register(r)
	user.Register(r)
}
