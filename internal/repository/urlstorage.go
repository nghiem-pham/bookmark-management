package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	defaultUrlExpTime = 24 * time.Hour
)

// UrlStorage is the interface for storing and retrieving shortened urls.
//
//go:generate mockery --name UrlStorage --filename urlstorage.go
type UrlStorage interface {
	StoreURL(ctx context.Context, code, url string, exp int) error
	GetURL(ctx context.Context, code string) (string, error)
}
type urlStorage struct {
	c *redis.Client
}

// NewUrlStorage creates a UrlStorage backed by the given redis.Client.
func NewUrlStorage(c *redis.Client) UrlStorage {
	return &urlStorage{c: c}
}

// StoreURL saves the mapping between code and url in Redis, setting it to
// expire after exp seconds. If exp is not positive, defaultUrlExpTime is
// used. It returns an error if the write fails.
func (s *urlStorage) StoreURL(ctx context.Context, code, url string, exp int) error {
	expiration := defaultUrlExpTime
	if exp > 0 {
		expiration = time.Duration(exp) * time.Second
	}
	return s.c.Set(ctx, code, url, expiration).Err()
}

// GetURL retrieves the url stored under code from Redis. It returns an
// error if the code does not exist or the read fails.
func (s *urlStorage) GetURL(ctx context.Context, code string) (string, error) {
	return s.c.Get(ctx, code).Result()
}
