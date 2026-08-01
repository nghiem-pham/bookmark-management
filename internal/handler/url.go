package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"

	"github.com/nghiem-pham/bookmark-management/internal/model"
	"github.com/nghiem-pham/bookmark-management/internal/service"
)

// UrlHandler is the interface for the url handler
type UrlHandler interface {
	Shorten(c *gin.Context)
	Redirect(c *gin.Context)
}

type urlHandler struct {
	shortenUrl service.ShortenUrl
}

// NewUrlHandler creates a new UrlHandler backed by the given
// service.ShortenUrl.
func NewUrlHandler(shortenUrl service.ShortenUrl) UrlHandler {
	return &urlHandler{
		shortenUrl: shortenUrl,
	}
}

// Shorten godoc
// @Summary Shorten a URL
// @Tags links
// @Accept json
// @Produce json
// @Param request body model.ShortenURLRequest true "URL to shorten"
// @Success 200 {object} model.ShortenURLResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/links/shorten [post]
func (h *urlHandler) Shorten(c *gin.Context) {
	var req model.ShortenURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	code, err := h.shortenUrl.ShortenURL(c.Request.Context(), req.URL, req.Exp)
	if err != nil {
		log.Error().Err(err).Msg("failed to shorten url")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, model.ShortenURLResponse{
		Code:    code,
		Message: "Shorten URL generated successfully!",
	})
}

// Redirect godoc
// @Summary Redirect to original URL
// @Tags links
// @Param code path string true "Shorten code"
// @Success 302
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/links/redirect/{code} [get]
func (h *urlHandler) Redirect(c *gin.Context) {
	code := c.Param("code")

	url, err := h.shortenUrl.GetURL(c.Request.Context(), code)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			c.JSON(http.StatusNotFound, gin.H{"error": "link not found"})
			return
		}

		log.Error().Err(err).Msg("failed to retrieve url")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.Redirect(http.StatusFound, url)
}
