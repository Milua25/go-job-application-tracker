package auth

import (
	"time"
)

// LoginRequest represents the request payload for user login.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest represents the request payload for user registration.
type RegisterRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8,max=64,strong_password"`
}

// RegisterResponse represents the response payload for user registration.
type RegisterResponse struct {
	ID        string    `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// LoginResponse represents the response payload for user login.
type LoginResponse struct {
	ID                    string    `json:"id"`
	FirstName             string    `json:"first_name"`
	LastName              string    `json:"last_name"`
	Email                 string    `json:"email"`
	SessionID             string    `json:"session_id"`
	AccessToken           string    `json:"access_token"`
	RefreshToken          string    `json:"refresh_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
}

// RenewAccessTokenReq represents the request payload for renewing an access token.
type RenewAccessTokenReq struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RenewAccessTokenResp represents the response payload for renewing an access token.
type RenewAccessTokenResp struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

// RefreshTokenRequest represents the request payload for refreshing an access token.
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// Actor represents the authenticated caller identity extracted at the HTTP boundary.
type Actor struct {
	UserID    string
	Email     string
	SessionID string
	IsAdmin   bool
}
