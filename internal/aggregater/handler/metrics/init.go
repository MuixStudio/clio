package metrics

import (
	"github.com/gin-gonic/gin"
	"github.com/muixstudio/clio/internal/aggregater/svc"
)

func RegisterWithoutAuth(r *gin.RouterGroup, svcCtx *svc.ServiceContext) {
	metricsHandler := NewMetricsHandler(svcCtx)
	r.GET("/metrics", metricsHandler.Metrics())
}
