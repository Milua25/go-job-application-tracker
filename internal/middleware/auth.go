package middleware

import (
	"errors"
	"log/slog"
	"strings"

	"github.com/Milua25/go-job-application-tracker/internal/authctx"
	"github.com/Milua25/go-job-application-tracker/internal/render"
	"github.com/Milua25/go-job-application-tracker/internal/token"
	"github.com/gin-gonic/gin"
)

// ErrTokenExpired is returned by TokenValidator when the token is past its expiry.
var (
	ErrTokenExpired                     = errors.New("token expired")
	ErrNilAuthService                   = errors.New("auth service is nil")
	ErrNoAuthorizationHeader            = errors.New("no Authorization header")
	ErrInvalidAuthorizationHeaderFormat = errors.New("invalid Authorization header format")
	ErrAuthTokenEmpty                   = errors.New("auth token is empty")
)

func AuthMiddleware(tokenMaker *token.JWTMaker) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := verifyClaimsFromAuthHeader(c, tokenMaker)
		if err != nil {
			if errors.Is(err, ErrTokenExpired) {
				render.UnauthorizedResponseError(c, "token has expired", token.ErrTokenExpired)
				return
			}
			if errors.Is(err, ErrNoAuthorizationHeader) {
				render.UnauthorizedResponseError(c, "authorization header is missing", ErrNoAuthorizationHeader)
				return
			}
			if errors.Is(err, ErrInvalidAuthorizationHeaderFormat) {
				render.UnauthorizedResponseError(c, "authorization header format must be 'Bearer <token>'", ErrInvalidAuthorizationHeaderFormat)
				return
			}
			if errors.Is(err, ErrAuthTokenEmpty) {
				render.UnauthorizedResponseError(c, "token is missing", ErrAuthTokenEmpty)
				return
			}
			if errors.Is(err, ErrNilAuthService) {
				render.InternalServerError(c, "auth service is nil", ErrNilAuthService)
				return
			}
			render.UnauthorizedResponseError(c, "invalid token", err)
			return
		}

		slog.Debug("token validated successfully", "user_id", claims.Uid, "email", claims.Email)
		c.Set(authctx.ContextKeyEmail, claims.Email)
		c.Set(authctx.ContextKeyFirstName, claims.FirstName)
		c.Set(authctx.ContextKeyLastName, claims.LastName)
		c.Set(authctx.ContextKeyUID, claims.Uid)
		c.Set(authctx.ContextKeySessionID, claims.SessionID)
		c.Set(authctx.ContextKeyIsAdmin, claims.IsAdmin)
		c.Next()
	}
}

func GetAdminMiddleware(tokenMaker *token.JWTMaker) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := verifyClaimsFromAuthHeader(c, tokenMaker)
		if err != nil {
			if errors.Is(err, ErrTokenExpired) {
				render.UnauthorizedResponseError(c, "token has expired", token.ErrTokenExpired)
				return
			}
			if errors.Is(err, ErrNoAuthorizationHeader) {
				render.UnauthorizedResponseError(c, "authorization header is missing", ErrNoAuthorizationHeader)
				return
			}
			if errors.Is(err, ErrInvalidAuthorizationHeaderFormat) {
				render.UnauthorizedResponseError(c, "authorization header format must be 'Bearer <token>'", ErrInvalidAuthorizationHeaderFormat)
				return
			}
			if errors.Is(err, ErrAuthTokenEmpty) {
				render.UnauthorizedResponseError(c, "token is missing", ErrAuthTokenEmpty)
				return
			}
			if errors.Is(err, ErrNilAuthService) {
				render.InternalServerError(c, "auth service is nil", ErrNilAuthService)
				return
			}
			render.UnauthorizedResponseError(c, "invalid token", err)
			return
		}

		if !claims.IsAdmin {
			slog.Warn("user is not an admin", "user_id", claims.Uid, "email", claims.Email)
			render.ForbiddenResponseError(c, "user is not an admin", errors.New("user is not an admin"))
			return
		}

		slog.Debug("token validated successfully", "user_id", claims.Uid, "email", claims.Email)
		c.Set(authctx.ContextKeyEmail, claims.Email)
		c.Set(authctx.ContextKeyFirstName, claims.FirstName)
		c.Set(authctx.ContextKeyLastName, claims.LastName)
		c.Set(authctx.ContextKeyUID, claims.Uid)
		c.Set(authctx.ContextKeySessionID, claims.SessionID)
		c.Set(authctx.ContextKeyIsAdmin, claims.IsAdmin)
		c.Next()
	}
}

func verifyClaimsFromAuthHeader(c *gin.Context, tokenMaker *token.JWTMaker) (*token.Claims, error) {
	slog.Debug("auth middleware invoked")
	if tokenMaker == nil {
		slog.Error("token maker is nil")
		return nil, ErrNilAuthService
	}

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		slog.Debug("missing authorization header")
		// render.UnauthorizedResponseError(c, "authorization header is missing", ErrNoAuthorizationHeader)
		return nil, ErrNoAuthorizationHeader
	}
	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		slog.Warn("invalid authorization header format")
		// render.UnauthorizedResponseError(c, "authorization header format must be 'Bearer <token>'", ErrInvalidAuthorizationHeaderFormat)
		return nil, ErrInvalidAuthorizationHeaderFormat
	}

	tokenString := authHeader[len(bearerPrefix):]
	if tokenString == "" {
		slog.Debug("empty token in authorization header")
		// render.UnauthorizedResponseError(c, "token is missing", ErrAuthTokenEmpty)
		return nil, ErrAuthTokenEmpty
	}

	claims, err := tokenMaker.ValidateToken(tokenString)
	if err != nil {
		if errors.Is(err, token.ErrTokenExpired) {
			slog.Debug("token expired")
			// render.UnauthorizedResponseError(c, "token has expired", token.ErrTokenExpired)
			return nil, ErrTokenExpired
		}
		slog.Warn("token validation failed", "error", err)
		// render.UnauthorizedResponseError(c, "invalid token", err)
		return nil, err
	}
	return claims, nil
}
