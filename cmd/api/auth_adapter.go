package main

import (
	"errors"

	"github.com/Milua25/go-job-application-tracker/internal/auth"
	"github.com/Milua25/go-job-application-tracker/internal/middleware"
)

// jwtAuthAdapter satisfies middleware.TokenValidator using tokens.AuthService.
// It lives here in the wiring layer so neither middleware nor tokens imports
// the other — the dependency graph stays acyclic by construction.
type jwtAuthAdapter struct {
	svc *auth.AuthService
}

func newJWTAuthAdapter(svc *auth.AuthService) middleware.TokenValidator {
	return &jwtAuthAdapter{svc: svc}
}

func (a *jwtAuthAdapter) ValidateToken(token string) (*middleware.Claims, error) {
	c, err := a.svc.ValidateToken(token)
	if err != nil {
		if errors.Is(err, auth.ErrTokenExpired) {
			return nil, middleware.ErrTokenExpired
		}
		return nil, err
	}
	return &middleware.Claims{
		Email:     c.Email,
		FirstName: c.FirstName,
		LastName:  c.LastName,
		Uid:       c.Uid,
	}, nil
}
