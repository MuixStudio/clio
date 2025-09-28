package health

import (
	"github.com/gin-gonic/gin"
	"github.com/muixstudio/clio/internal/aggregater/svc"
)

func RegisterWithoutAuth(r *gin.RouterGroup, svcCtx *svc.ServiceContext) {

	authHandler := NewHealthCheckHandler(svcCtx)

	r.GET("/health_check", authHandler.HealthCheck())
}
