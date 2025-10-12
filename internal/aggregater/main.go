package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/metrics"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/muixstudio/clio/internal/aggregater/handler"
	"github.com/muixstudio/clio/internal/aggregater/middleware/logger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

var (
	// Name is the name of the compiled software.
	Name = "metrics"
	// Version is the version of the compiled software.
	// Version = "v1.0.0"

	_metricRequests metric.Int64Counter
	_metricSeconds  metric.Float64Histogram
)

func main() {

	exporter, err := prometheus.New()
	if err != nil {
		panic(err)
	}
	provider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(exporter))
	meter := provider.Meter(Name)

	_metricRequests, err = metrics.DefaultRequestsCounter(meter, metrics.DefaultServerRequestsCounterName)
	if err != nil {
		panic(err)
	}

	_metricSeconds, err = metrics.DefaultSecondsHistogram(meter, metrics.DefaultServerSecondsHistogramName)
	if err != nil {
		panic(err)
	}

	r := gin.New()
	r.Use(logger.Logger())

	gin.DisableBindValidation()
	gin.SetMode(gin.ReleaseMode)

	handler.Register(&r.RouterGroup)

	httpSrv := http.NewServer(
		http.Address(":5020"),
		http.Middleware(
			metrics.Server(
				metrics.WithSeconds(_metricSeconds),
				metrics.WithRequests(_metricRequests),
			),
		),
	)
	httpSrv.Handle("/metrics", promhttp.Handler())
	httpSrv.HandlePrefix("/", r)

	app := kratos.New(
		kratos.Name("aggregater"),
		kratos.Server(
			httpSrv,
		),
	)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
