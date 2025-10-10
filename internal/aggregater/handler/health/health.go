package health

import (
	"github.com/gin-gonic/gin"
	"github.com/muixstudio/clio/internal/aggregater/svc"
	"github.com/muixstudio/clio/internal/aggregater/utils/response"
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
		response.SuccessOK(c)
	}
}
