package service

import (
	"context"

	"github.com/nghiem-pham/bookmark-management/internal/repository"
	"github.com/nghiem-pham/bookmark-management/pkg/stringutils"
)

const (
	urlCodeLength = 7
)

// ShortenUrl is the interface for the url shortening service
//
//go:generate mockery --name ShortenUrl --filename urlstorage.go
type ShortenUrl interface {
	ShortenURL(ctx context.Context, url string) (string, error)
}
type shortenUrl struct {
	repo repository.UrlStorage
}

func NewShortenUrl(repo repository.UrlStorage) ShortenUrl {
	return &shortenUrl{repo: repo}
}

func (s *shortenUrl) ShortenURL(ctx context.Context, url string) (string, error) {
	// generate key
	urlCode, err := stringutils.GenerateCode(urlCodeLength)
	if err != nil {
		return "", err
	}

	// store in repo
	err = s.repo.StoreURL(ctx, urlCode, url)
	if err != nil {
		return "", err
	}

	return urlCode, nil
}
