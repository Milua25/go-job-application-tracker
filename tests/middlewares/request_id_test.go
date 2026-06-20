package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Milua25/go-job-application-tracker/internal/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newRequestIDTestRouter() *gin.Engine {
	router := gin.New()
	router.Use(middlewares.RequestIDMiddleware())
	router.GET("/health", func(c *gin.Context) {
		id, _ := c.Get("request_id")
		c.JSON(http.StatusOK, gin.H{"request_id": id})
	})
	return router
}

func TestRequestID_GeneratedWhenAbsent(t *testing.T) {
	router := newRequestIDTestRouter()

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/health", nil)
	require.NoError(t, err)
	router.ServeHTTP(w, req)

	id := w.Header().Get("X-Request-ID")
	assert.NotEmpty(t, id)
	_, err = uuid.Parse(id)
	assert.NoError(t, err, "generated X-Request-ID should be a valid UUID")
}

func TestRequestID_EchoedWhenPresent(t *testing.T) {
	router := newRequestIDTestRouter()
	incoming := "my-custom-request-id"

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/health", nil)
	require.NoError(t, err)
	req.Header.Set("X-Request-ID", incoming)
	router.ServeHTTP(w, req)

	assert.Equal(t, incoming, w.Header().Get("X-Request-ID"))
}

func TestRequestID_StoredInContext(t *testing.T) {
	router := newRequestIDTestRouter()
	incoming := "ctx-check-id"

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/health", nil)
	require.NoError(t, err)
	req.Header.Set("X-Request-ID", incoming)
	router.ServeHTTP(w, req)

	assert.Contains(t, w.Body.String(), incoming)
}
