package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Milua25/go-job-application-tracker/internal/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func newCorsTestRouter() *gin.Engine {
	router := gin.New()
	router.Use(middlewares.Cors())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	return router
}

func TestCors_AllowedOrigin(t *testing.T) {
	router := newCorsTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/health", nil)
	req.Header.Set("Origin", "http://localhost:8000")
	router.ServeHTTP(w, req)

	if w.Header().Get("Access-Control-Allow-Origin") != "http://localhost:8000" {
		t.Fatalf("expected Access-Control-Allow-Origin to be http://localhost:8000, got %s", w.Header().Get("Access-Control-Allow-Origin"))
	}
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 OK, got %d", w.Code)
	}
}

func TestCors_AllowedMethods(t *testing.T) {
	router := newCorsTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodOptions, "/health", nil)
	req.Header.Set("Origin", "http://localhost:8000")
	req.Header.Set("Access-Control-Request-Method", "GET")
	router.ServeHTTP(w, req)

	assert.Equal(t, "GET,POST,PUT,DELETE", w.Header().Get("Access-Control-Allow-Methods"))
	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestCors_AllowHeaders(t *testing.T) {
	router := newCorsTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodOptions, "/health", nil)
	req.Header.Set("Origin", "http://localhost:8000")
	req.Header.Set("Access-Control-Request-Headers", "X-CSRF-TOKEN")
	router.ServeHTTP(w, req)

	assert.Equal(t, "Origin,Content-Type,Authorization,X-Csrf-Token", w.Header().Get("Access-Control-Allow-Headers"))
	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestCors_Preflight(t *testing.T) {
	router := newCorsTestRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodOptions, "/health", nil)
	req.Header.Set("Origin", "http://localhost:8000")
	req.Header.Set("Access-Control-Request-Method", "POST")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}
	if w.Header().Get("Access-Control-Allow-Origin") != "http://localhost:8000" {
		t.Fatalf("missing Allow-Origin on preflight")
	}
}
func TestCors_DisallowedOrigin(t *testing.T) {
	router := newCorsTestRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/health", nil)
	req.Header.Set("Origin", "http://evil.com")
	router.ServeHTTP(w, req)

	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("expected no Access-Control-Allow-Origin, got %s", got)
	}
}

func TestCors_AllowCredentials(t *testing.T) {
	router := newCorsTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/health", nil)
	req.Header.Set("Origin", "http://localhost:8000")
	router.ServeHTTP(w, req)

	assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
}

func newRequireOriginTestRouter() *gin.Engine {
	router := gin.New()
	router.Use(middlewares.RequireOrigin())
	router.OPTIONS("/health", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})
	router.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	return router
}

func TestRequireOrigin_OPTIONS_WithoutOrigin(t *testing.T) {
	router := newRequireOriginTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodOptions, "/health", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRequireOrigin_OPTIONS_WithOrigin(t *testing.T) {
	router := newRequireOriginTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodOptions, "/health", nil)
	req.Header.Set("Origin", "http://localhost:8000")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestRequireOrigin_NonOPTIONS_WithoutOrigin(t *testing.T) {
	router := newRequireOriginTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/health", nil)
	router.ServeHTTP(w, req)

	// GET without Origin is not blocked — only OPTIONS is checked
	assert.Equal(t, http.StatusOK, w.Code)
}
