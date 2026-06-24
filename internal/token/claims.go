package token

import (
	"time"

	"github.com/Milua25/go-job-application-tracker/internal/user"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Email     string
	FirstName string
	LastName  string
	Uid       string
	IsAdmin   bool
	SessionID string
	jwt.RegisteredClaims
}

func newUserClaims(user *user.User, sessionID, issuer string, expireDuration time.Duration) *Claims {
	now := time.Now()
	return &Claims{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Uid:       user.ID.String(),
		IsAdmin:   user.IsAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{"access"},
			ExpiresAt: jwt.NewNumericDate(now.Add(expireDuration)),
			Issuer:    issuer,
			IssuedAt:  jwt.NewNumericDate(now),
		},
		SessionID: sessionID,
	}
}

func newRefreshClaims(user *user.User, issuer string, expireDuration time.Duration) *Claims {
	now := time.Now()
	return &Claims{
		Uid: user.ID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{"refresh"},
			ExpiresAt: jwt.NewNumericDate(now.Add(expireDuration)),
			Issuer:    issuer,
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
}
