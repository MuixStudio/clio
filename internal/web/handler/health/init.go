package health

import (
	"github.com/gin-gonic/gin"
	"github.com/muixstudio/clio/internal/web/svc"
)

func RegisterWithoutAuth(r *gin.RouterGroup, svcCtx *svc.ServiceContext) {

	healthCheckHandler := NewHealthCheckHandler(svcCtx)

	r.GET("/health_check", healthCheckHandler.HealthCheck())
}
