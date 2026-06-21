package middleware

import (
	"errors"
	"strings"

	"github.com/Milua25/go-job-application-tracker/internal/render"
	"github.com/gin-gonic/gin"
)

// Context keys set by AuthMiddleware. Use these constants instead of bare
// strings to avoid silent collisions with other middleware.
const (
	ContextKeyUID       = "uid"
	ContextKeyEmail     = "email"
	ContextKeyFirstName = "first_name"
	ContextKeyLastName  = "last_name"
)

// Claims carries the verified identity fields AuthMiddleware needs from a token.
type Claims struct {
	Email     string
	FirstName string
	LastName  string
	Uid       string
}

// ErrTokenExpired is returned by TokenValidator when the token is past its expiry.
var (
	ErrTokenExpired                     = errors.New("token expired")
	ErrNilAuthService                   = errors.New("auth service is nil")
	ErrNoAuthorizationHeader            = errors.New("no Authorization header")
	ErrInvalidAuthorizationHeaderFormat = errors.New("invalid Authorization header format")
	ErrAuthTokenEmpty                   = errors.New("auth token is empty")
)

// TokenValidator is the only contract AuthMiddleware depends on.
// Concrete implementations live in cmd/api so neither this package
// nor the auth package imports the other.
type TokenValidator interface {
	ValidateToken(token string) (*Claims, error)
}

func AuthMiddleware(authService TokenValidator) gin.HandlerFunc {
	return func(c *gin.Context) {
		if authService == nil {
			render.InternalServerError(c, "auth service is not initialized", ErrNilAuthService)
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			render.UnAuthorizedResponseError(c, "authorization header is missing", ErrNoAuthorizationHeader)
			return
		}
		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			render.UnAuthorizedResponseError(c, "authorization header format must be 'Bearer <token>'", ErrInvalidAuthorizationHeaderFormat)
			return
		}

		tokenString := authHeader[len(bearerPrefix):]
		if tokenString == "" {
			render.UnAuthorizedResponseError(c, "token is missing", ErrAuthTokenEmpty)
			return
		}
		claims, err := authService.ValidateToken(tokenString)
		if err != nil {
			if errors.Is(err, ErrTokenExpired) {
				render.UnAuthorizedResponseError(c, "token has expired", ErrTokenExpired)
				return
			}
			render.UnAuthorizedResponseError(c, "invalid token", err)
			return
		}

		c.Set(ContextKeyEmail, claims.Email)
		c.Set(ContextKeyFirstName, claims.FirstName)
		c.Set(ContextKeyLastName, claims.LastName)
		c.Set(ContextKeyUID, claims.Uid)
		c.Next()
	}
}
