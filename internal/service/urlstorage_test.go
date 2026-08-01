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
		exp int

		setupMockRepo func(t *testing.T) *mocks.UrlStorage

		expectedErr bool
	}{
		{
			name: "success",
			url:  "https://example.com",
			exp:  604800,

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
			exp:  604800,

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

			code, err := svc.ShortenURL(ctx, tc.url, tc.exp)

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

func TestShortenUrl_GetURL(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	testCases := []struct {
		name string

		code string

		setupMockRepo func(t *testing.T) *mocks.UrlStorage

		expectedURL string
		expectedErr bool
	}{
		{
			name: "success",
			code: "1234567",

			setupMockRepo: func(t *testing.T) *mocks.UrlStorage {
				repoMock := mocks.NewUrlStorage(t)
				repoMock.On("GetURL", ctx, "1234567").
					Return("https://example.com", nil)
				return repoMock
			},

			expectedURL: "https://example.com",
			expectedErr: false,
		},
		{
			name: "repository error",
			code: "notexist",

			setupMockRepo: func(t *testing.T) *mocks.UrlStorage {
				repoMock := mocks.NewUrlStorage(t)
				repoMock.On("GetURL", ctx, "notexist").
					Return("", errors.New("redis: nil"))
				return repoMock
			},

			expectedURL: "",
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repoMock := tc.setupMockRepo(t)
			svc := NewShortenUrl(repoMock)

			url, err := svc.GetURL(ctx, tc.code)

			if tc.expectedErr {
				assert.Error(t, err)
				assert.Empty(t, url)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedURL, url)
			}
		})
	}
}
