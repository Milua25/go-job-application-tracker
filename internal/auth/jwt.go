package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/Milua25/go-job-application-tracker/internal/user"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Email     string
	FirstName string
	LastName  string
	Uid       string
	jwt.RegisteredClaims
}

type AuthService struct {
	secretKey       string
	issuer          string
	expireIn        string
	refreshExpireIn string
}

// Define custom errors for token validation
var (
	ErrTokenExpired           = errors.New("token has expired")
	ErrInvalidToken           = errors.New("invalid token")
	ErrTokenEmpty             = errors.New("token is empty")
	ErrValidation             = errors.New("error validating token")
	ErrTokenMalformed         = errors.New("malformed token")
	ErrIncorrectSigningMethod = errors.New("incorrect signing method")
)

func newAuthService(secretKey, issuer, expireIn, refreshExpireIn string) *AuthService {
	return &AuthService{
		secretKey:       secretKey,
		issuer:          issuer,
		expireIn:        expireIn,
		refreshExpireIn: refreshExpireIn,
	}
}

// GenerateToken generates a JWT token with the provided user information and expiration times.
func (s *AuthService) GenerateToken(user *user.User) (string, string, error) {
	expireDuration, err := time.ParseDuration(s.expireIn)
	if err != nil {
		return "", "", err
	}

	refreshExpireDuration, err := time.ParseDuration(s.refreshExpireIn)
	if err != nil {
		return "", "", err
	}
	now := time.Now()
	tokenClaims := Claims{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Uid:       user.ID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{"access"},
			ExpiresAt: jwt.NewNumericDate(now.Add(expireDuration)),
			Issuer:    s.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	refreshTokenClaims := Claims{
		Uid: user.ID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{"refresh"},
			ExpiresAt: jwt.NewNumericDate(now.Add(refreshExpireDuration)),
			Issuer:    s.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	// Create a new token object, specifying signing method and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)

	// Sign and get the complete encoded token as a string using the secret
	signedToken, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", "", fmt.Errorf("error signing token: %w", err)
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)

	signedRefreshToken, err := refreshToken.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", "", fmt.Errorf("error signing refresh token: %w", err)
	}
	return signedToken, signedRefreshToken, nil

}

// ValidateToken validates an access token and returns the claims if valid.
func (s *AuthService) ValidateToken(signedToken string) (*Claims, error) {
	if signedToken == "" {
		return nil, ErrTokenEmpty
	}
	// Guard against an empty secret or issuer set at construction time.
	if s.secretKey == "" || s.issuer == "" {
		return nil, fmt.Errorf("auth service is not properly initialized")
	}

	return validateToken(signedToken, s.secretKey, s.issuer, "access")
}

// ValidateRefreshToken validates a refresh token and returns the claims if valid.
func (s *AuthService) ValidateRefreshToken(signedToken string) (*Claims, error) {
	if signedToken == "" {
		return nil, ErrTokenEmpty
	}
	if s.secretKey == "" || s.issuer == "" {
		return nil, fmt.Errorf("auth service is not properly initialized")
	}

	return validateToken(signedToken, s.secretKey, s.issuer, "refresh")
}

func validateToken(signedToken, secretkey, issuer, audience string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(signedToken, &Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrIncorrectSigningMethod
		}
		return []byte(secretkey), nil
	}, jwt.WithIssuer(issuer), jwt.WithAudience(audience), jwt.WithExpirationRequired(), jwt.WithIssuedAt())

	if token == nil {
		return nil, ErrInvalidToken
	}

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, ErrTokenMalformed
		}
		return nil, ErrValidation
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	if claims.Uid == "" {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
