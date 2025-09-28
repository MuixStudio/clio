package health

import (
	"github.com/gin-gonic/gin"
	"github.com/muixstudio/clio/internal/aggregater/svc"
)

type HealthCheckHandler struct {
	svcCtx *svc.ServiceContext
}

func NewHealthCheckHandler(svcCtx *svc.ServiceContext) *HealthCheckHandler {
	return &HealthCheckHandler{
		svcCtx: svcCtx,
	}
}

func (ah HealthCheckHandler) HealthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}
