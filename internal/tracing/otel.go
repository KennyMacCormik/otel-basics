package tracing

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/trace"
)

type Config struct {
	Endpoint string `mapstructure:"trace_endpoint" validate:"url" env:"TRACE_ENDPOINT"`
}

func NewTraceProvider(ctx context.Context, cfg Config, regGlobal bool) (*trace.TracerProvider, error) {
	te, err := otlptracehttp.New(ctx, otlptracehttp.WithEndpointURL(cfg.Endpoint))
	if err != nil {
		return nil, err
	}
	tp := trace.NewTracerProvider(trace.WithBatcher(te))
	if regGlobal {
		otel.SetTracerProvider(tp)
	}
	return tp, nil
}
