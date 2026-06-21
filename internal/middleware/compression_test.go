package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Milua25/go-job-application-tracker/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func newCompressionTestRouter() *gin.Engine {
	router := gin.New()
	router.Use(middleware.Compression())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	router.GET("/api/v1/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	return router
}

func TestCompression_WithGzip(t *testing.T) {
	router := newCompressionTestRouter()

	// Test with gzip encoding
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/health", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	router.ServeHTTP(w, req)

	// Check if the response is gzipped
	assert.Equal(t, "gzip", w.Header().Get("Content-Encoding"))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCompression_ExcludedPath(t *testing.T) {
	router := newCompressionTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/health", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Empty(t, w.Header().Get("Content-Encoding"), "excluded path should not be gzip compressed")
}

// client sends no Accept-Encoding → response should be plain
func TestCompression_WithoutGzip(t *testing.T) {
	r := newCompressionTestRouter()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if enc := w.Header().Get("Content-Encoding"); enc != "" {
		t.Fatalf("expected no Content-Encoding, got %s", enc)
	}
	if w.Body.String() != `{"status":"ok"}` {
		t.Fatalf("unexpected body: %s", w.Body.String())
	}
}
