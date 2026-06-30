package token

import (
	"time"

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

func newUserClaims(userID, email, firstName, lastName string, isAdmin bool, sessionID, issuer string, expireDuration time.Duration) *Claims {
	now := time.Now()
	return &Claims{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Uid:       userID,
		IsAdmin:   isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{"access"},
			ExpiresAt: jwt.NewNumericDate(now.Add(expireDuration)),
			Issuer:    issuer,
			IssuedAt:  jwt.NewNumericDate(now),
		},
		SessionID: sessionID,
	}
}

func newRefreshClaims(userID, issuer string, expireDuration time.Duration) *Claims {
	now := time.Now()
	return &Claims{
		Uid: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{"refresh"},
			ExpiresAt: jwt.NewNumericDate(now.Add(expireDuration)),
			Issuer:    issuer,
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
}
