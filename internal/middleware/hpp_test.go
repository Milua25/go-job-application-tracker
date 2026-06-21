package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Milua25/go-job-application-tracker/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newHPPTestRouter(opts middleware.HPPOptions) *gin.Engine {
	router := gin.New()
	router.Use(opts.Hpp())

	router.GET("/search", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"raw_query": c.Request.URL.RawQuery})
	})
	router.POST("/submit", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		c.JSON(http.StatusOK, gin.H{
			"foo": c.Request.PostForm["foo"],
			"bar": c.Request.PostForm["bar"],
		})
	})

	return router
}

func TestHPP_QueryDeduplication(t *testing.T) {
	opts := middleware.HPPOptions{CheckQuery: true, Whitelist: []string{"foo"}}
	router := newHPPTestRouter(opts)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/search?foo=a&foo=b", nil)
	require.NoError(t, err)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	// duplicate foo values deduplicated to the first
	assert.Contains(t, w.Body.String(), "foo=a")
	assert.NotContains(t, w.Body.String(), "foo=b")
}

func TestHPP_QueryWhitelistFiltering(t *testing.T) {
	opts := middleware.HPPOptions{CheckQuery: true, Whitelist: []string{"foo"}}
	router := newHPPTestRouter(opts)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/search?foo=a&bar=x", nil)
	require.NoError(t, err)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	// bar is not whitelisted and should be removed from the query
	assert.Contains(t, w.Body.String(), "foo=a")
	assert.NotContains(t, w.Body.String(), "bar=x")
}

func TestHPP_QueryFilterDisabled(t *testing.T) {
	opts := middleware.HPPOptions{CheckQuery: false, Whitelist: []string{"foo"}}
	router := newHPPTestRouter(opts)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/search?foo=a&foo=b&bar=x", nil)
	require.NoError(t, err)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	// no filtering — all params preserved
	assert.Contains(t, w.Body.String(), "foo=a")
	assert.Contains(t, w.Body.String(), "bar=x")
}

func TestHPP_BodyDeduplication(t *testing.T) {
	opts := middleware.HPPOptions{
		CheckBody:        true,
		CheckContentType: "application/x-www-form-urlencoded",
		Whitelist:        []string{"foo"},
	}
	router := newHPPTestRouter(opts)

	w := httptest.NewRecorder()
	body := strings.NewReader("foo=a&foo=b")
	req, err := http.NewRequest(http.MethodPost, "/submit", body)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	// duplicate foo values deduplicated to the first
	assert.Contains(t, w.Body.String(), `"foo":["a"]`)
}

func TestHPP_BodySkippedOnGET(t *testing.T) {
	opts := middleware.HPPOptions{
		CheckBody:        true,
		CheckContentType: "application/x-www-form-urlencoded",
		Whitelist:        []string{"foo"},
	}
	router := newHPPTestRouter(opts)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/search?foo=a&foo=b", nil)
	require.NoError(t, err)
	router.ServeHTTP(w, req)

	// CheckBody only applies to POST/PUT/PATCH/DELETE; GET should pass through unchanged
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHPP_BodySkippedOnWrongContentType(t *testing.T) {
	opts := middleware.HPPOptions{
		CheckBody:        true,
		CheckContentType: "application/x-www-form-urlencoded",
		Whitelist:        []string{"foo"},
	}
	router := newHPPTestRouter(opts)

	w := httptest.NewRecorder()
	body := strings.NewReader(`{"foo":["a","b"]}`)
	req, err := http.NewRequest(http.MethodPost, "/submit", body)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// body filtering is skipped when content type does not match
	assert.Equal(t, http.StatusOK, w.Code)
}
