package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nghiem-pham/bookmark-management/internal/service"
	"github.com/nghiem-pham/bookmark-management/internal/service/mocks"
	"github.com/stretchr/testify/assert"
)

var testErr = errors.New("test error")

func TestHealthHandler(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupRequest     func(ctx *gin.Context)
		setupMockService func(ctx context.Context) *mocks.HealthService

		expectedStatus   int
		expectedResponse string
	}{
		{
			name: "success",

			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/health-check", nil)
			},

			setupMockService: func(ctx context.Context) *mocks.HealthService {
				serviceMock := mocks.NewHealthService(t)
				serviceMock.On("HealthCheck").Return(&service.HealthCheckResponse{
					Message:     "OK",
					ServiceName: "bookmark_service",
					InstanceID:  "test-instance-id",
				})
				return serviceMock
			},

			expectedStatus:   http.StatusOK,
			expectedResponse: `{"instance_id":"test-instance-id","message":"OK","service_name":"bookmark_service"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)
			tc.setupRequest(ctx)

			mockSvc := tc.setupMockService(ctx)
			testHandler := NewHealthHandler(mockSvc)

			testHandler.HealthCheck(ctx)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.Equal(t, tc.expectedResponse, rec.Body.String())
		})
	}
}
