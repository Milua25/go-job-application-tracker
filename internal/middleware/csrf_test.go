package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Milua25/go-job-application-tracker/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newCSRFTestRouter() *gin.Engine {
	router := gin.New()
	router.Use(middleware.SessionCSRFProtection("test-session-secret"))
	router.Use(middleware.CSRFProtection("test-csrf-secret"))
	router.Use(middleware.CSRFTokenHeader())

	router.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	router.POST("/resource", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	return router
}

// fetchCSRFToken performs a GET to obtain a valid CSRF token and its session cookie.
func fetchCSRFToken(t *testing.T, router *gin.Engine) (token string, cookies []*http.Cookie) {
	t.Helper()
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/health", nil)
	require.NoError(t, err)
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	token = w.Header().Get("X-CSRF-TOKEN")
	require.NotEmpty(t, token, "GET should return a CSRF token")
	return token, w.Result().Cookies()
}

func TestCSRFTokenHeader_SetOnGET(t *testing.T) {
	router := newCSRFTestRouter()
	token, _ := fetchCSRFToken(t, router)
	assert.NotEmpty(t, token)
}

func TestCSRFTokenHeader_NotSetOnPOST(t *testing.T) {
	router := newCSRFTestRouter()
	token, cookies := fetchCSRFToken(t, router)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/resource", nil)
	require.NoError(t, err)
	req.Header.Set("X-CSRF-TOKEN", token)
	for _, c := range cookies {
		req.AddCookie(c)
	}
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Empty(t, w.Header().Get("X-CSRF-TOKEN"))
}

func TestCSRFProtection_GETPassthrough(t *testing.T) {
	router := newCSRFTestRouter()

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/health", nil)
	require.NoError(t, err)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCSRFProtection_ValidToken(t *testing.T) {
	router := newCSRFTestRouter()
	token, cookies := fetchCSRFToken(t, router)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/resource", nil)
	require.NoError(t, err)
	req.Header.Set("X-CSRF-TOKEN", token)
	for _, c := range cookies {
		req.AddCookie(c)
	}
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCSRFProtection_MissingToken(t *testing.T) {
	router := newCSRFTestRouter()

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/resource", nil)
	require.NoError(t, err)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCSRFProtection_InvalidToken(t *testing.T) {
	router := newCSRFTestRouter()
	_, cookies := fetchCSRFToken(t, router)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/resource", nil)
	require.NoError(t, err)
	req.Header.Set("X-CSRF-TOKEN", "not-a-valid-token")
	for _, c := range cookies {
		req.AddCookie(c)
	}
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCSRFProtection_TokenWithoutSession(t *testing.T) {
	router := newCSRFTestRouter()
	token, _ := fetchCSRFToken(t, router)

	// Send a valid token but no session cookie — should be rejected.
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/resource", nil)
	require.NoError(t, err)
	req.Header.Set("X-CSRF-TOKEN", token)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
