package main

import (
	"context"

	"github.com/gin-gonic/gin"
	kzap "github.com/go-kratos/kratos/contrib/log/zap/v2"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/metrics"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/muixstudio/clio/internal/aggregater/handler"
	"github.com/muixstudio/clio/internal/aggregater/middleware/logger"
	ginMiddleware "github.com/muixstudio/clio/internal/aggregater/middleware/metrics"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	// Name is the name of the compiled software.
	Name = "metrics"
	// Version is the version of the compiled software.
	// Version = "v1.0.0"

	_metricRequests metric.Int64Counter
	_metricSeconds  metric.Float64Histogram
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

	exporter, err := initMetricsExporter(context.Background())
	if err != nil {
		panic(err)
	}
	provider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter)))
	meter := provider.Meter(Name)

	_metricRequests, err = metrics.DefaultRequestsCounter(meter, metrics.DefaultServerRequestsCounterName)
	if err != nil {
		panic(err)
	}

	_metricSeconds, err = metrics.DefaultSecondsHistogram(meter, metrics.DefaultServerSecondsHistogramName)
	if err != nil {
		panic(err)
	}

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
	)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
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

func initMetricsExporter(ctx context.Context) (sdkmetric.Exporter, error) {
	// 配置 OTLP endpoint
	// 默认端点: localhost:4317 (gRPC)
	// 如果使用 HTTP: localhost:4318

	conn, err := grpc.NewClient(
		"localhost:4317", // OTLP gRPC endpoint
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	exporter, err := otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithGRPCConn(conn),
	)
	if err != nil {
		return nil, err
	}

	return exporter, nil
}
