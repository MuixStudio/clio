package observability

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

var res = resource.NewWithAttributes(
	semconv.SchemaURL,
	semconv.ServiceName("aggregater-service"), // 服务名
	semconv.ServiceVersion("1.2.3"),           // 版本
	semconv.ServiceInstanceID("instance-001"), // 实例ID
	semconv.ServiceNamespace("clio"),          // 环境
	attribute.String("cluster", "us-west-2"),  // 集群
	attribute.String("team", "backend"),       // 团队
)
