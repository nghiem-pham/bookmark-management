package repository

import (
	"context"
	"testing"
	"time"

	redisPkg "github.com/nghiem-pham/bookmark-management/pkg/redis"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestUrlStorage_StoreURL(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		exp int

		setupMock func() *redis.Client

		expectErr  error
		verifyFunc func(ctx context.Context, r *redis.Client)
	}{
		{
			name: "normal case",
			exp:  0,

			setupMock: func() *redis.Client {
				mock := redisPkg.InitMockRedis(t)
				return mock
			},

			expectErr: nil,
			verifyFunc: func(ctx context.Context, r *redis.Client) {
				url, err := r.Get(ctx, "1234567").Result()
				assert.Nil(t, err)
				assert.Equal(t, url, "https://google.com")
			},
		},
		{
			name: "with expiration",
			exp:  604800,

			setupMock: func() *redis.Client {
				return redisPkg.InitMockRedis(t)
			},

			expectErr: nil,
			verifyFunc: func(ctx context.Context, r *redis.Client) {
				ttl, err := r.TTL(ctx, "1234567").Result()
				assert.Nil(t, err)
				assert.Equal(t, 604800*time.Second, ttl)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			redisMock := tc.setupMock()
			testRepo := NewUrlStorage(redisMock)

			err := testRepo.StoreURL(ctx, "1234567", "https://google.com", tc.exp)
			assert.Equal(t, tc.expectErr, err)
			if err == nil {
				tc.verifyFunc(ctx, redisMock)
			}
		})
	}
}

func TestUrlStorage_GetURL(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		code string

		setupMock func() *redis.Client

		expectURL string
		expectErr error
	}{
		{
			name: "normal case",
			code: "1234567",

			setupMock: func() *redis.Client {
				mock := redisPkg.InitMockRedis(t)
				err := mock.Set(context.Background(), "1234567", "https://google.com", 0).Err()
				assert.Nil(t, err)
				return mock
			},

			expectURL: "https://google.com",
			expectErr: nil,
		},
		{
			name: "code not found",
			code: "notexist",
			setupMock: func() *redis.Client {
				return redisPkg.InitMockRedis(t)
			},

			expectURL: "",
			expectErr: redis.Nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			redisMock := tc.setupMock()
			testRepo := NewUrlStorage(redisMock)

			url, err := testRepo.GetURL(ctx, tc.code)
			assert.Equal(t, tc.expectErr, err)
			assert.Equal(t, tc.expectURL, url)
		})
	}

}
