package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/nghiem-pham/bookmark-management/internal/handler"
	"github.com/nghiem-pham/bookmark-management/internal/service"
)

type Engine interface {
	Start() error
	ServeHTTP(w *httptest.ResponseRecorder, req *http.Request)
}

type engine struct {
	app *gin.Engine
	cfg *Config
}

func NewEngine(cfg *Config) Engine {
	e := &engine{
		app: gin.Default(),
		cfg: cfg,
	}
	e.initRoutes()

	return e
}

// Start starts the application
func (e *engine) Start() error {
	return e.app.Run(fmt.Sprintf(":%s", e.cfg.AppPort))
}

// ServeHTTP to test the API endpoint
func (e *engine) ServeHTTP(w *httptest.ResponseRecorder, req *http.Request) {
	e.app.ServeHTTP(w, req)
}

// initRoutes initializes the routes
func (e *engine) initRoutes() {
	healthSvc := service.NewHealthService(e.cfg.ServiceName, e.cfg.InstanceID)
	healthHandler := handler.NewHealthHandler(healthSvc)

	e.app.GET("/health-check", healthHandler.HealthCheck)
}
