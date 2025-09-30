package exporters

import (
	"os"

	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
)

func ExporterDevNull() (*tracesdk.TracerProvider, error) {
	sink, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return nil, err
	}
	exp, err := stdouttrace.New(stdouttrace.WithWriter(sink))
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
		)),
		tracesdk.WithBatcher(exp),
	)

	return tp, nil
}
