package exporters

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
)

func ExporterJaegerTCP(url, serviceName, environment string, attrs ...attribute.KeyValue) (*tracesdk.TracerProvider, error) {
	attrs = append(attrs, semconv.ServiceName(serviceName), attribute.String(envKey, environment))

	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewWithAttributes(semconv.SchemaURL, attrs...)))

	return tp, nil
}
