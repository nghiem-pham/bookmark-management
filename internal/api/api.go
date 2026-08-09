package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/nghiem-pham/bookmark-management/internal/handler"
	"github.com/nghiem-pham/bookmark-management/internal/repository"
	"github.com/nghiem-pham/bookmark-management/internal/service"
	redisPkg "github.com/nghiem-pham/bookmark-management/pkg/redis"
	"github.com/redis/go-redis/v9"
)

// Engine defines the interface for the application server, exposing
// methods to start the server and handle HTTP requests.
type Engine interface {
	Start() error
	ServeHTTP(w *httptest.ResponseRecorder, req *http.Request)
}

type engine struct {
	app *gin.Engine
	cfg *Config
}

// NewEngine creates a new Engine instance configured with the given Config
// and registers all application routes.
func NewEngine(cfg *Config) (Engine, error) {
	redisClient, err := redisPkg.NewClient("URLSTORAGE")
	if err != nil {
		return nil, err
	}

	e := &engine{
		app: gin.Default(),
		cfg: cfg,
	}
	e.initRoutes(redisClient)

	return e, nil
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
func (e *engine) initRoutes(redisClient *redis.Client) {
	healthSvc := service.NewHealthService(e.cfg.ServiceName, e.cfg.InstanceID)
	healthHandler := handler.NewHealthHandler(healthSvc)
	e.app.GET("/health-check", healthHandler.HealthCheck)

	urlStorage := repository.NewUrlStorage(redisClient)
	shortenUrl := service.NewShortenUrl(urlStorage)
	e.registerUrlRoutes(shortenUrl)
}

// initRoutesWithServices initializes routes with injected services
func (e *engine) initRoutesWithServices(shortenUrl service.ShortenUrl) {
	healthSvc := service.NewHealthService(e.cfg.ServiceName, e.cfg.InstanceID)
	healthHandler := handler.NewHealthHandler(healthSvc)
	e.app.GET("/health-check", healthHandler.HealthCheck)

	e.registerUrlRoutes(shortenUrl)
}

// registerUrlRoutes registers the url shortening routes
func (e *engine) registerUrlRoutes(shortenUrl service.ShortenUrl) {
	urlHandler := handler.NewUrlHandler(shortenUrl)
	e.app.POST("/v1/links/shorten", urlHandler.Shorten)
	e.app.GET("/v1/links/redirect/:code", urlHandler.Redirect)
}
