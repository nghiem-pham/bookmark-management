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
	ShortenURL(ctx context.Context, url string, exp int) (string, error)
	GetURL(ctx context.Context, code string) (string, error)
}
type shortenUrl struct {
	repo repository.UrlStorage
}

// NewShortenUrl creates a ShortenUrl service backed by the given
// repository.UrlStorage.
func NewShortenUrl(repo repository.UrlStorage) ShortenUrl {
	return &shortenUrl{repo: repo}
}

// ShortenURL generates a random code of length urlCodeLength for url,
// stores the code-to-url mapping in the repository, and returns the
// generated code. It returns an error if code generation or storage fails.
func (s *shortenUrl) ShortenURL(ctx context.Context, url string, exp int) (string, error) {
	// generate key
	urlCode, err := stringutils.GenerateCode(urlCodeLength)
	if err != nil {
		return "", err
	}

	// store in repo
	err = s.repo.StoreURL(ctx, urlCode, url, exp)
	if err != nil {
		return "", err
	}

	return urlCode, nil
}

// GetURL retrieves the original url for code from the repository. It
// returns an error if the code does not exist or the read fails.
func (s *shortenUrl) GetURL(ctx context.Context, code string) (string, error) {
	return s.repo.GetURL(ctx, code)
}
