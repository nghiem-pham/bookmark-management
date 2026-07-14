package integration_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nghiem-pham/bookmark-management/internal/api"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestHealthCheckEndpoint(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		serviceName string
		instanceID  string
	}{
		{
			name:        "health check returns correct values",
			serviceName: "bookmark_service",
			instanceID:  "integration-test-instance",
		},
		{
			name:        "health check with custom config",
			serviceName: "custom_service",
			instanceID:  "custom-instance-id",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			cfg := &api.Config{
				ServiceName: tc.serviceName,
				InstanceID:  tc.instanceID,
				AppPort:     "8080",
			}
			app := api.NewEngine(cfg)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/health-check", nil)
			app.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			var resp map[string]string
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err)
			assert.Equal(t, "OK", resp["message"])
			assert.Equal(t, tc.serviceName, resp["service_name"])
			assert.Equal(t, tc.instanceID, resp["instance_id"])
		})
	}
}
