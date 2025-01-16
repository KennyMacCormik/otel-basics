package grpc

import (
	"context"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"jaeger/internal/db/net/grpc/StorageEndpoint"
	"jaeger/internal/db/repo"
	"jaeger/internal/proto/db"
	"log/slog"
	"strings"
)

const tracerName = "db"

func NewGrpcServer(st repo.Storage, lg *slog.Logger) *grpc.Server {
	gs := grpc.NewServer(grpc.UnaryInterceptor(tracing()),
		grpc.StatsHandler(otelgrpc.NewServerHandler()))
	db.RegisterStorageEndpointServer(gs, StorageEndpoint.NewStorageEndpoint(st, lg))
	return gs
}

func tracing() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}

		// Normalize incoming metadata keys to lowercase
		lowercaseMD := metadata.New(nil)
		for k, v := range md {
			lowercaseMD[strings.ToLower(k)] = v
		}

		propagator := otel.GetTextMapPropagator()
		ctx = propagator.Extract(ctx, propagation.HeaderCarrier(md))

		tracer := otel.Tracer(tracerName)
		ctx, span := tracer.Start(ctx, "db"+info.FullMethod)
		defer span.End()

		span.AddEvent("test event")

		return handler(ctx, req)
	}
}
