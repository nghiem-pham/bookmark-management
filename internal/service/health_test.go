package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthService_HealthCheck(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		serviceName string
		instanceID  string
	}{
		{
			name:        "success",
			serviceName: "bookmark_service",
			instanceID:  "test-instance-id",
		},
		{
			name:        "custom service name",
			serviceName: "custom_service",
			instanceID:  "another-instance",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := NewHealthService(tc.serviceName, tc.instanceID)
			resp := svc.HealthCheck()

			assert.Equal(t, "OK", resp.Message)
			assert.Equal(t, tc.serviceName, resp.ServiceName)
			assert.Equal(t, tc.instanceID, resp.InstanceID)
		})
	}
}
