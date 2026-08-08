package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"

	"github.com/nghiem-pham/bookmark-management/internal/model"
	"github.com/nghiem-pham/bookmark-management/internal/service/mocks"
)

func TestShorten(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupRequest     func(ctx *gin.Context)
		setupMockService func(ctx *gin.Context) *mocks.ShortenUrl

		expectedStatus   int
		expectedResponse string
	}{
		{
			name: "success",

			setupRequest: func(ctx *gin.Context) {
				body, _ := json.Marshal(model.ShortenURLRequest{
					URL: "https://example.com/long/url",
					Exp: 604800,
				})
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/links/shorten", bytes.NewBuffer(body))
			},

			setupMockService: func(ctx *gin.Context) *mocks.ShortenUrl {
				serviceMock := mocks.NewShortenUrl(t)
				serviceMock.On("ShortenURL", ctx, "https://example.com/long/url", 604800).
					Return("abc1234", nil)
				return serviceMock
			},

			expectedStatus:   http.StatusOK,
			expectedResponse: `{"code":"abc1234","message":"Shorten URL generated successfully!"}`,
		},
		{
			name: "invalid request body",

			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/links/shorten", bytes.NewBuffer([]byte(`{"url":""}`)))
			},

			setupMockService: func(ctx *gin.Context) *mocks.ShortenUrl {
				return mocks.NewShortenUrl(t)
			},

			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "internal server error",

			setupRequest: func(ctx *gin.Context) {
				body, _ := json.Marshal(model.ShortenURLRequest{
					URL: "https://example.com/long/url",
					Exp: 604800,
				})
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/links/shorten", bytes.NewBuffer(body))
			},

			setupMockService: func(ctx *gin.Context) *mocks.ShortenUrl {
				serviceMock := mocks.NewShortenUrl(t)
				serviceMock.On("ShortenURL", ctx, "https://example.com/long/url", 604800).
					Return("", assert.AnError)
				return serviceMock
			},

			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)
			tc.setupRequest(ctx)

			mockSvc := tc.setupMockService(ctx)
			testHandler := NewUrlHandler(mockSvc)

			testHandler.Shorten(ctx)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			if tc.expectedResponse != "" {
				assert.Equal(t, tc.expectedResponse, rec.Body.String())
			}
		})
	}
}

func TestRedirect(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupRequest     func(ctx *gin.Context)
		setupMockService func(ctx *gin.Context) *mocks.ShortenUrl

		expectedStatus int
		expectedURL    string
	}{
		{
			name: "success",

			setupRequest: func(ctx *gin.Context) {
				ctx.Params = gin.Params{{Key: "code", Value: "abc1234"}}
				ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/links/redirect/abc1234", nil)
			},

			setupMockService: func(ctx *gin.Context) *mocks.ShortenUrl {
				serviceMock := mocks.NewShortenUrl(t)
				serviceMock.On("GetURL", ctx, "abc1234").
					Return("https://example.com/long/url", nil)
				return serviceMock
			},

			expectedStatus: http.StatusFound,
			expectedURL:    "https://example.com/long/url",
		},
		{
			name: "link not found",

			setupRequest: func(ctx *gin.Context) {
				ctx.Params = gin.Params{{Key: "code", Value: "notfound"}}
				ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/links/redirect/notfound", nil)
			},

			setupMockService: func(ctx *gin.Context) *mocks.ShortenUrl {
				serviceMock := mocks.NewShortenUrl(t)
				serviceMock.On("GetURL", ctx, "notfound").
					Return("", redis.Nil)
				return serviceMock
			},

			expectedStatus: http.StatusNotFound,
		},
		{
			name: "internal server error",

			setupRequest: func(ctx *gin.Context) {
				ctx.Params = gin.Params{{Key: "code", Value: "abc1234"}}
				ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/links/redirect/abc1234", nil)
			},

			setupMockService: func(ctx *gin.Context) *mocks.ShortenUrl {
				serviceMock := mocks.NewShortenUrl(t)
				serviceMock.On("GetURL", ctx, "abc1234").
					Return("", assert.AnError)
				return serviceMock
			},

			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)
			tc.setupRequest(ctx)

			mockSvc := tc.setupMockService(ctx)
			testHandler := NewUrlHandler(mockSvc)

			testHandler.Redirect(ctx)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			if tc.expectedURL != "" {
				assert.Equal(t, tc.expectedURL, rec.Header().Get("Location"))
			}
		})
	}
}
