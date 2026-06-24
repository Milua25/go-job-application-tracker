package token

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/Milua25/go-job-application-tracker/internal/user"

	"github.com/golang-jwt/jwt/v5"
)

type JWTMaker struct {
	secretKey       string
	issuer          string
	expireDuration  time.Duration
	refreshDuration time.Duration
}

// NewJWTMaker creates a new JWTMaker. Returns an error if either duration string is invalid.
func NewJWTMaker(secretKey, issuer, expireIn, refreshExpireIn string) (*JWTMaker, error) {
	expireDuration, err := time.ParseDuration(expireIn)
	if err != nil {
		return nil, fmt.Errorf("invalid expireIn %q: %w", expireIn, err)
	}
	refreshDuration, err := time.ParseDuration(refreshExpireIn)
	if err != nil {
		return nil, fmt.Errorf("invalid refreshExpireIn %q: %w", refreshExpireIn, err)
	}
	return &JWTMaker{
		secretKey:       secretKey,
		issuer:          issuer,
		expireDuration:  expireDuration,
		refreshDuration: refreshDuration,
	}, nil
}

// GenerateToken generates an access token for the given user.
func (s *JWTMaker) GenerateToken(u *user.User, sessionID string) (string, time.Time, error) {
	tokenClaims := newUserClaims(u, sessionID, s.issuer, s.expireDuration)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)
	signedToken, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("error signing token: %w", err)
	}

	return signedToken, tokenClaims.ExpiresAt.Time, nil
}

// GenerateRefreshToken generates a refresh token for the given user.
func (s *JWTMaker) CreateRefreshToken(u *user.User) (string, time.Time, time.Time, error) {
	refreshTokenClaims := newRefreshClaims(u, s.issuer, s.refreshDuration)

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	signedRefreshToken, err := refreshToken.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", time.Time{}, time.Time{}, fmt.Errorf("error signing refresh token: %w", err)
	}
	// get the expiration time of the refresh token
	refreshTokenClaims, ok := refreshToken.Claims.(*Claims)
	if !ok {
		return "", time.Time{}, time.Time{}, fmt.Errorf("error extracting claims from refresh token")
	}
	slog.Debug("refresh token generated successfully", "user_id", u.ID.String(), "expires_at", refreshTokenClaims.ExpiresAt.Time)

	return signedRefreshToken, refreshTokenClaims.IssuedAt.Time, refreshTokenClaims.ExpiresAt.Time, nil
}

// ValidateToken validates an access token and returns the claims if valid.
func (s *JWTMaker) ValidateToken(signedToken string) (*Claims, error) {
	slog.Debug("validating access token")
	if signedToken == "" {
		slog.Warn("empty access token provided")
		return nil, ErrTokenEmpty
	}
	if s.secretKey == "" || s.issuer == "" {
		slog.Error("token maker not properly initialized", "has_secret", s.secretKey != "", "has_issuer", s.issuer != "")
		return nil, ErrServiceNotInitialized
	}
	return validateToken(signedToken, s.secretKey, s.issuer, "access")
}

// ValidateRefreshToken validates a refresh token and returns the claims if valid.
func (s *JWTMaker) ValidateRefreshToken(signedToken string) (*Claims, error) {
	slog.Debug("validating refresh token")
	if signedToken == "" {
		slog.Warn("empty refresh token provided")
		return nil, ErrTokenEmpty
	}
	if s.secretKey == "" || s.issuer == "" {
		slog.Error("token maker not properly initialized", "has_secret", s.secretKey != "", "has_issuer", s.issuer != "")
		return nil, ErrServiceNotInitialized
	}
	return validateToken(signedToken, s.secretKey, s.issuer, "refresh")
}

func validateToken(signedToken, secretKey, issuer, audience string) (*Claims, error) {
	slog.Debug("parsing and validating token", "audience", audience)

	parsedToken, err := jwt.ParseWithClaims(
		signedToken,
		&Claims{},
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				slog.Warn("incorrect signing method", "method", token.Method.Alg())
				return nil, ErrIncorrectSigningMethod
			}
			return []byte(secretKey), nil
		}, jwt.WithIssuer(issuer), jwt.WithAudience(audience), jwt.WithExpirationRequired(), jwt.WithIssuedAt())

	if parsedToken == nil {
		slog.Warn("token parsing resulted in nil token", "audience", audience)
		return nil, ErrInvalidToken
	}

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			slog.Debug("token has expired", "audience", audience)
			return nil, ErrTokenExpired
		}
		if errors.Is(err, jwt.ErrTokenMalformed) {
			slog.Warn("token is malformed", "audience", audience)
			return nil, ErrTokenMalformed
		}
		slog.Warn("token validation failed", "audience", audience, "error", err)
		return nil, ErrValidation
	}

	claims, ok := parsedToken.Claims.(*Claims)
	if !ok || !parsedToken.Valid {
		slog.Warn("failed to extract or invalid claims", "audience", audience)
		return nil, ErrInvalidToken
	}

	if claims.Uid == "" {
		slog.Warn("token claims missing user id", "audience", audience)
		return nil, ErrInvalidToken
	}

	slog.Debug("token validated successfully", "audience", audience, "user_id", claims.Uid)
	return claims, nil
}
