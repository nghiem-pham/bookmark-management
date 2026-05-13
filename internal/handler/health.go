package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nghiem-pham/bookmark-management/internal/service"
)

// HealthHandler is the interface for the health handler
type HealthHandler interface {
	HealthCheck(c *gin.Context)
}

type healthHandler struct {
	healthService service.HealthService
}

func NewHealthHandler(healthSvc service.HealthService) HealthHandler {
	return &healthHandler{
		healthService: healthSvc,
	}
}

// HealthCheck godoc
// @Summary Health check
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health-check [get]
func (h *healthHandler) HealthCheck(c *gin.Context) {
	resp := h.healthService.HealthCheck()
	c.JSON(http.StatusOK, gin.H{
		"message":      resp.Message,
		"service_name": resp.ServiceName,
		"instance_id":  resp.InstanceID,
	})
}
