package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/muixstudio/clio/internal/aggregater/handler/auth"
	"github.com/muixstudio/clio/internal/aggregater/handler/health"
	"github.com/muixstudio/clio/internal/aggregater/handler/metrics"
	"github.com/muixstudio/clio/internal/aggregater/handler/user"
	"github.com/muixstudio/clio/internal/aggregater/middleware"
	"github.com/muixstudio/clio/internal/aggregater/middleware/cors"
	"github.com/muixstudio/clio/internal/aggregater/svc"
)

func Register(r *gin.RouterGroup) {

	svcCtx := svc.NewServiceContext()

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
	metrics.RegisterWithoutAuth(webCallWithoutAuth, svcCtx)

	webCall := r.Group("")
	webCall.Use(middleware.WebCallAuth())
	auth.Register(webCall, svcCtx)
	user.Register(webCall, svcCtx)

	// api
	api := r.Group("/api")
	user.Register(api, svcCtx)

}
