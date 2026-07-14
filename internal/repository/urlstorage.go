package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	urlExpTime = 24 * time.Hour
)

// UrlStorage is the interface for storing and retrieving shortened urls.
//
//go:generate mockery --name UrlStorage --filename urlstorage.go
type UrlStorage interface {
	StoreURL(ctx context.Context, code, url string) error
	GetURL(ctx context.Context, code string) (string, error)
}
type urlStorage struct {
	c *redis.Client
}

func NewUrlStorage(c *redis.Client) UrlStorage {
	return &urlStorage{c: c}
}

func (s *urlStorage) StoreURL(ctx context.Context, code, url string) error {
	return s.c.Set(ctx, code, url, urlExpTime).Err()
}

func (s *urlStorage) GetURL(ctx context.Context, code string) (string, error) {
	return s.c.Get(ctx, code).Result()
}
