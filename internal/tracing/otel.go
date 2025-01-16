package tracing

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
)

const serviceName = "jaeger-demo-app"

type Config struct {
	Endpoint string `mapstructure:"trace_endpoint" validate:"url" env:"TRACE_ENDPOINT"`
}

func NewTraceProvider(ctx context.Context, cfg Config, regGlobal bool) (*trace.TracerProvider, error) {
	te, err := otlptracehttp.New(ctx, otlptracehttp.WithEndpointURL(cfg.Endpoint))
	if err != nil {
		return nil, err
	}
	tp := trace.NewTracerProvider(
		trace.WithBatcher(te),
		trace.WithResource(resource.NewSchemaless(
			semconv.ServiceNameKey.String(serviceName),
		)),
	)
	if regGlobal {
		otel.SetTracerProvider(tp)
	}
	otel.SetTextMapPropagator(propagation.TraceContext{})
	return tp, nil
}
