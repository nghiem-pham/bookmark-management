package service

import (
	"context"
	"errors"
	"testing"

	"github.com/nghiem-pham/bookmark-management/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestShortenUrl_ShortenURL(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	testCases := []struct {
		name string

		url string

		setupMockRepo func(t *testing.T) *mocks.UrlStorage

		expectedErr bool
	}{
		{
			name: "success",
			url:  "https://example.com",

			setupMockRepo: func(t *testing.T) *mocks.UrlStorage {
				repoMock := mocks.NewUrlStorage(t)
				repoMock.On("StoreURL", ctx, mock.AnythingOfType("string"), "https://example.com").
					Return(nil)
				return repoMock
			},

			expectedErr: false,
		},
		{
			name: "repository error",
			url:  "https://example.com",

			setupMockRepo: func(t *testing.T) *mocks.UrlStorage {
				repoMock := mocks.NewUrlStorage(t)
				repoMock.On("StoreURL", ctx, mock.AnythingOfType("string"), "https://example.com").
					Return(errors.New("redis error"))
				return repoMock
			},

			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repoMock := tc.setupMockRepo(t)
			svc := NewShortenUrl(repoMock)

			code, err := svc.ShortenURL(ctx, tc.url)

			if tc.expectedErr {
				assert.Error(t, err)
				assert.Empty(t, code)
			} else {
				assert.NoError(t, err)
				assert.Len(t, code, urlCodeLength)
			}
		})
	}
}
