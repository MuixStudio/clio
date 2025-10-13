package metrics

import (
	"github.com/gin-gonic/gin"
	"github.com/muixstudio/clio/internal/aggregater/svc"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type MetricsHandler struct {
	svcCtx *svc.ServiceContext
}

func NewMetricsHandler(svcCtx *svc.ServiceContext) *MetricsHandler {
	return &MetricsHandler{
		svcCtx: svcCtx,
	}
}

func (ah MetricsHandler) Metrics() gin.HandlerFunc {
	return gin.WrapH(promhttp.Handler())
}
