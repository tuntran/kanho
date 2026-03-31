package storage

import (
	"context"
	"io"
	"time"
)

type Client interface {
	PutObject(ctx context.Context, key string, reader io.Reader, size int64, contentType string) error
	GetPresignedURL(ctx context.Context, key string, expires time.Duration) (string, error)
	GetPresignedPutURL(ctx context.Context, key string, expires time.Duration) (string, error)
	DeleteObject(ctx context.Context, key string) error
	Ping(ctx context.Context) error
}
