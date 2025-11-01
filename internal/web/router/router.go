package router

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/muixstudio/clio/internal/web/handler/auth"
	"github.com/muixstudio/clio/internal/web/handler/health"
	"github.com/muixstudio/clio/internal/web/handler/user"
	"github.com/muixstudio/clio/internal/web/middleware"
	"github.com/muixstudio/clio/internal/web/middleware/cors"
	"github.com/muixstudio/clio/internal/web/middleware/logger"
	ginMiddleware "github.com/muixstudio/clio/internal/web/middleware/metrics"
	"github.com/muixstudio/clio/internal/web/svc"
)

func initEngine() *gin.Engine {
	r := gin.New()
	r.Use(
		logger.Logger(),
		ginMiddleware.Metrics("clio"),
	)

	gin.DisableBindValidation()
	gin.SetMode(gin.ReleaseMode)
	return r
}

func NewEngine(ctx context.Context, svcCtx *svc.ServiceContext) *gin.Engine {
	r := initEngine()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type"},
		//ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// webcall
	webCallWithoutAuth := r.Group("")
	auth.RegisterWithoutAuth(webCallWithoutAuth, svcCtx)
	health.RegisterWithoutAuth(webCallWithoutAuth, svcCtx)

	webCall := r.Group("")
	webCall.Use(middleware.WebCallAuth())
	auth.Register(webCall, svcCtx)
	user.Register(webCall, svcCtx)

	// api
	api := r.Group("/api")
	user.Register(api, svcCtx)
	return r
}
