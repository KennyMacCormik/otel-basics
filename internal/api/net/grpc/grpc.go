package grpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"jaeger/internal/proto/db"
	"time"
)

type StorageEndpointClient interface {
	Get(key string) (string, bool, error)
	Set(key, value string) error
	Del(key string) error
}
type Client struct {
	client                *grpc.ClientConn
	storageEndpointClient db.StorageEndpointClient
	timeout               time.Duration
}

func NewGrpcClient(addr string, t time.Duration) (*Client, error) {
	c, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return &Client{}, fmt.Errorf("grpc client error: %w", err)
	}
	return &Client{client: c, storageEndpointClient: db.NewStorageEndpointClient(c), timeout: t}, nil
}

func (c *Client) Close() error {
	return c.client.Close()
}

func (c *Client) Get(key string) (string, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	kv, err := c.storageEndpointClient.Get(ctx, &db.Key{Key: key})
	if err != nil {
		return "", false, err
	}
	return kv.GetVal(), true, nil
}

func (c *Client) Set(key, value string) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	_, err := c.storageEndpointClient.Set(ctx, &db.KeyValue{Key: key, Val: value})
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Del(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	_, err := c.storageEndpointClient.Del(ctx, &db.Key{Key: key})
	if err != nil {
		return err
	}
	return nil
}
