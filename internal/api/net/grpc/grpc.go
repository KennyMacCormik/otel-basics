package grpc

import (
	"context"
	"fmt"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"jaeger/internal/proto/db"
	"strings"
	"time"
)

type StorageEndpointClient interface {
	Get(ctx context.Context, key string) (string, bool, error)
	Set(ctx context.Context, key, value string) error
	Del(ctx context.Context, key string) error
}
type Client struct {
	client                *grpc.ClientConn
	storageEndpointClient db.StorageEndpointClient
	timeout               time.Duration
}

func NewGrpcClient(addr string, t time.Duration) (*Client, error) {
	c, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		return &Client{}, fmt.Errorf("grpc client error: %w", err)
	}
	return &Client{client: c, storageEndpointClient: db.NewStorageEndpointClient(c), timeout: t}, nil
}

func (c *Client) Close() error {
	return c.client.Close()
}

func (c *Client) Get(ctx context.Context, key string) (string, bool, error) {
	grpcCtx, cancel := prepGrpcCtx(ctx, c.timeout)
	defer cancel()

	kv, err := c.storageEndpointClient.Get(grpcCtx, &db.Key{Key: key})
	if err != nil {
		return "", false, err
	}
	return kv.GetVal(), true, nil
}

func (c *Client) Set(ctx context.Context, key, value string) error {
	grpcCtx, cancel := prepGrpcCtx(ctx, c.timeout)
	defer cancel()

	if _, err := c.storageEndpointClient.Set(grpcCtx, &db.KeyValue{Key: key, Val: value}); err != nil {
		return err
	}

	return nil
}

func (c *Client) Del(ctx context.Context, key string) error {
	grpcCtx, cancel := prepGrpcCtx(ctx, c.timeout)
	defer cancel()

	if _, err := c.storageEndpointClient.Del(grpcCtx, &db.Key{Key: key}); err != nil {
		return err
	}

	return nil
}

func normalizeMetadata(md metadata.MD) metadata.MD {
	lowercaseMD := metadata.New(nil)
	for k, v := range md {
		lowercaseMD[strings.ToLower(k)] = v
	}
	return lowercaseMD
}

func prepGrpcCtx(ctx context.Context, t time.Duration) (context.Context, context.CancelFunc) {
	md := metadata.New(nil)
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(md))

	md = normalizeMetadata(md)

	tctx := metadata.NewOutgoingContext(ctx, md)
	return context.WithTimeout(tctx, t)
}
