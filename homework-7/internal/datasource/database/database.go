package database

import (
	"context"
	"time"
)

type Client struct {
}

func (c *Client) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	return nil
}

func (c *Client) Get(ctx context.Context, key string) (any, error) {
	return nil, nil
}
