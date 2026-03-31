package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

type LocalClient struct {
	basePath string
}

func NewLocalClient(basePath string) (*LocalClient, error) {
	if err := os.MkdirAll(basePath, 0o755); err != nil {
		return nil, fmt.Errorf("creating storage dir: %w", err)
	}
	return &LocalClient{basePath: basePath}, nil
}

func (l *LocalClient) PutObject(_ context.Context, key string, reader io.Reader, _ int64, _ string) error {
	fullPath := filepath.Join(l.basePath, key)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		return err
	}
	f, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, reader)
	return err
}

func (l *LocalClient) GetPresignedURL(_ context.Context, key string, _ time.Duration) (string, error) {
	return "/files/" + key, nil
}

func (l *LocalClient) GetPresignedPutURL(_ context.Context, key string, _ time.Duration) (string, error) {
	return "/files/upload/" + key, nil
}

func (l *LocalClient) DeleteObject(_ context.Context, key string) error {
	return os.Remove(filepath.Join(l.basePath, key))
}

func (l *LocalClient) Ping(_ context.Context) error {
	_, err := os.Stat(l.basePath)
	return err
}
