package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Milua25/go-job-application-tracker/internal/authctx"
	"github.com/Milua25/go-job-application-tracker/internal/middleware"
	"github.com/Milua25/go-job-application-tracker/internal/token"
	"github.com/Milua25/go-job-application-tracker/internal/user"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newAuthTokenMaker(t *testing.T, expireIn string) *token.JWTMaker {
	t.Helper()

	maker, err := token.NewJWTMaker("test-secret-key-1234567890", "test-issuer", expireIn, "24h")
	require.NoError(t, err)
	return maker
}

func newAuthTestUser(isAdmin bool) *user.User {
	return &user.User{
		ID:        uuid.New(),
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
		IsAdmin:   isAdmin,
	}
}

func issueAccessToken(t *testing.T, maker *token.JWTMaker, isAdmin bool) string {
	t.Helper()

	tok, _, err := maker.GenerateToken(newAuthTestUser(isAdmin), uuid.NewString())
	require.NoError(t, err)
	return tok
}

func newAuthRouter(tokenMaker *token.JWTMaker) *gin.Engine {
	router := gin.New()
	router.Use(middleware.AuthMiddleware(tokenMaker))
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"uid":        c.GetString(authctx.ContextKeyUID),
			"email":      c.GetString(authctx.ContextKeyEmail),
			"session_id": c.GetString(authctx.ContextKeySessionID),
		})
	})
	return router
}

func newAdminRouter(tokenMaker *token.JWTMaker) *gin.Engine {
	router := gin.New()
	router.Use(middleware.GetAdminMiddleware(tokenMaker))
	router.GET("/admin", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	return router
}

func TestAuthMiddleware_MissingAuthorizationHeader(t *testing.T) {
	router := newAuthRouter(newAuthTokenMaker(t, "15m"))

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/protected", nil)
	require.NoError(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_InvalidAuthorizationHeaderFormat(t *testing.T) {
	router := newAuthRouter(newAuthTokenMaker(t, "15m"))

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/protected", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Token abc")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_EmptyBearerToken(t *testing.T) {
	router := newAuthRouter(newAuthTokenMaker(t, "15m"))

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/protected", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer ")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_ExpiredToken(t *testing.T) {
	maker := newAuthTokenMaker(t, "-1m")
	expiredToken := issueAccessToken(t, maker, false)
	router := newAuthRouter(maker)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/protected", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+expiredToken)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_NilTokenMakerReturnsInternalServerError(t *testing.T) {
	router := newAuthRouter(nil)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/protected", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer any-token")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetAdminMiddleware_NonAdminReturnsForbidden(t *testing.T) {
	maker := newAuthTokenMaker(t, "15m")
	tok := issueAccessToken(t, maker, false)
	router := newAdminRouter(maker)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/admin", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+tok)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestGetAdminMiddleware_AdminPasses(t *testing.T) {
	maker := newAuthTokenMaker(t, "15m")
	tok := issueAccessToken(t, maker, true)
	router := newAdminRouter(maker)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/admin", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+tok)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
