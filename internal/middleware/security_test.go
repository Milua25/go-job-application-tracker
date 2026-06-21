package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Milua25/go-job-application-tracker/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func newSecurityTestRouter() *gin.Engine {
	router := gin.New()
	router.Use(middleware.SecurityHeaders())
	router.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	return router
}

func TestSecurityHeaders(t *testing.T) {
	router := newSecurityTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/health", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	headers := []struct {
		name     string
		expected string
	}{
		{"X-Frame-Options", "DENY"},
		{"X-XSS-Protection", "1; mode=block"},
		{"Content-Security-Policy", "default-src 'self'"},
		{"Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload"},
		{"Referrer-Policy", "no-referrer"},
		{"Permissions-Policy", "geolocation=(), microphone=()"},
		{"Cross-Origin-Opener-Policy", "same-origin"},
		{"Cross-Origin-Embedder-Policy", "require-corp"},
		{"Cross-Origin-Resource-Policy", "same-origin"},
		{"Cache-Control", "no-store"},
		{"X-Content-Type-Options", "nosniff"},
	}

	for _, h := range headers {
		t.Run(h.name, func(t *testing.T) {
			assert.Equal(t, h.expected, w.Header().Get(h.name))
		})
	}
}
