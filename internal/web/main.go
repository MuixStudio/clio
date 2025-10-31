package main

import (
	"context"

	"github.com/gin-gonic/gin"
	kzap "github.com/go-kratos/kratos/contrib/log/zap/v2"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/muixstudio/clio/internal/web/handler"
	"github.com/muixstudio/clio/internal/web/middleware/logger"
	ginMiddleware "github.com/muixstudio/clio/internal/web/middleware/metrics"
	mm "github.com/muixstudio/clio/internal/web/svc/observability"
	"go.uber.org/zap"
)

type testWriteSyncer struct {
	output []string
}

func (x *testWriteSyncer) Write(p []byte) (n int, err error) {
	x.output = append(x.output, string(p))
	return len(p), nil
}

func (x *testWriteSyncer) Sync() error {
	return nil
}

func main() {
	mm.InitMeterProvider()
	httpSrv := http.NewServer(
		http.Address(":5020"),
	)

	r := initGinRouter()
	httpSrv.HandlePrefix("/", r)

	//core := zapcore.NewCore(zapcore.NewJSONEncoder(zapcore.EncoderConfig{
	//	MessageKey:     "msg",
	//	LevelKey:       "level",
	//	NameKey:        "logger",
	//	EncodeLevel:    zapcore.LowercaseLevelEncoder,
	//	EncodeTime:     zapcore.ISO8601TimeEncoder,
	//	EncodeDuration: zapcore.StringDurationEncoder,
	//}), &testWriteSyncer{}, zap.DebugLevel)
	zapLogger := zap.NewExample()
	kLogger := kzap.NewLogger(zapLogger)

	app := kratos.New(
		kratos.Name("aggregater"),
		kratos.Version("v1.0.0"),
		kratos.Logger(kLogger),
		kratos.Server(
			httpSrv,
		),
		kratos.BeforeStart(beforeStart),
	)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

func beforeStart(ctx context.Context) error {
	return nil
}

func initGinRouter() *gin.Engine {
	r := gin.New()
	r.Use(
		logger.Logger(),
		ginMiddleware.Metrics("clio"),
	)

	gin.DisableBindValidation()
	gin.SetMode(gin.ReleaseMode)

	handler.Register(&r.RouterGroup)
	return r
}
