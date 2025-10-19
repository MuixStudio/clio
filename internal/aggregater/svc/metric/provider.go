package metric

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	// Name is the name of the compiled software.
	Name = "metrics"
	// Version is the version of the compiled software.
	// Version = "v1.0.0"
)

func initProvider() {
	exporter, err := initMetricsExporter(context.Background())
	if err != nil {
		panic(err)
	}
	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter)),
		sdkmetric.WithResource(resource.NewWithAttributes("https://opentelemetry.io/schemas/1.11.0", attribute.String("service.name1", "aggregater"))),
	)
	_ = provider.Meter(
		Name,
		metric.WithInstrumentationVersion("v1.0.0"),
		metric.WithSchemaURL("https://opentelemetry.io/schemas/1.11.0"),
		metric.WithInstrumentationAttributes(
			attribute.String("job", "aggregater"),
			attribute.String("instance", "aggregater.instance"),
		),
	)
	otel.SetMeterProvider(provider)
}

func initMetricsExporter(ctx context.Context) (sdkmetric.Exporter, error) {
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
