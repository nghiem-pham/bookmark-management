package service

// HealthCheckResponse is the response for the health check endpoint
type HealthCheckResponse struct {
	Message     string
	ServiceName string
	InstanceID  string
}

// HealthService is the interface for the health service
//
//go:generate mockery --name HealthService --filename health.go
type HealthService interface {
	HealthCheck() *HealthCheckResponse
}

type healthService struct {
	serviceName string
	instanceID  string
}

func NewHealthService(serviceName, instanceID string) HealthService {
	return &healthService{
		serviceName: serviceName,
		instanceID:  instanceID,
	}
}

func (s *healthService) HealthCheck() *HealthCheckResponse {
	return &HealthCheckResponse{
		Message:     "OK",
		ServiceName: s.serviceName,
		InstanceID:  s.instanceID,
	}
}
