package main

import (
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/gin-gonic/gin"
	kitexzap "github.com/kitex-contrib/obs-opentelemetry/logging/zap"
	"github.com/muixstudio/clio/internal/aggregater/handler"
	"github.com/muixstudio/clio/internal/aggregater/middleware/logger"
)

func main() {
	klog.SetLevel(klog.LevelDebug)
	zapLogger := kitexzap.NewLogger(kitexzap.WithCustomFields("app", "aggregater"))
	klog.SetLogger(zapLogger)
	r := gin.New()
	r.Use(logger.Logger())

	gin.DisableBindValidation()
	gin.SetMode(gin.ReleaseMode)

	handler.Register(&r.RouterGroup)

	if err := r.Run(":5020"); err != nil {
		panic(err)
	}
}
