package metric

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

type AppMetrics struct {
	requestCounter metric.Int64Counter
}

func NewAppMetrics() *AppMetrics {
	initProvider()
	meter := otel.Meter(Name)

	requestCounter, _ := meter.Int64Counter(
		"http_server_requests",
		metric.WithDescription("Total number of HTTP requests"),
		metric.WithUnit("{request}"),
	)

	requestCounter.Add(context.Background(), 8)

	return &AppMetrics{
		requestCounter: requestCounter,
	}
}
