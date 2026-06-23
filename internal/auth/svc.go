package auth

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/Milua25/go-job-application-tracker/internal/token"
	"github.com/Milua25/go-job-application-tracker/internal/user"
	"github.com/Milua25/go-job-application-tracker/pkg/utils"
	"github.com/google/uuid"
)

type authService struct {
	userStore    user.Repository
	sessionStore token.Repository
	tokenMaker   *token.JWTMaker
}

func newAuthService(store user.Repository, sessionStore token.Repository, tokenMaker *token.JWTMaker) *authService {
	return &authService{
		userStore:    store,
		sessionStore: sessionStore,
		tokenMaker:   tokenMaker,
	}
}

func (h *authService) register(ctx context.Context, req RegisterRequest) (*user.User, error) {

	existingUser, err := h.userStore.GetByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, ErrNotFound) {
		slog.Error("failed to check existing user", "error", err)
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	newUser := user.User{
		ID:           uuid.New(),
		Email:        req.Email,
		PasswordHash: hashedPassword,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Timezone:     "UTC",
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err = h.userStore.Create(ctx, &newUser); err != nil {
		return nil, err
	}

	return &newUser, nil
}

func (h *authService) loginUser(ctx context.Context, req LoginRequest) (*user.User, string, string, error) {

	foundUser, err := h.userStore.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			return nil, "", "", ErrInvalidCredentials
		}
		slog.Error("failed to retrieve user by email", "email", req.Email, "error", err)
		return nil, "", "", err
	}

	// checks if the provided password matches the stored password hash
	if !utils.VerifyPassword(foundUser.PasswordHash, req.Password) {
		slog.Warn("authentication failed", "email", req.Email)
		return nil, "", "", ErrInvalidCredentials
	}

	accessToken, err := h.tokenMaker.GenerateToken(foundUser)
	if err != nil {
		slog.Error("failed to generate token", "email", foundUser.Email, "error", err)
		return nil, "", "", err
	}

	// delete any existing session before creating a new one to prevent unique constraint violations on refresh_token
	err = h.sessionStore.DeleteSessionsByEmail(ctx, foundUser.Email)
	if err != nil && !errors.Is(err, token.ErrNotFound) {
		slog.Error("failed to delete existing session", "email", foundUser.Email, "error", err)
		return nil, "", "", ErrFailedToDeleteSession
	}

	refreshToken, refreshTokenIssuedAt, refreshTokenExpiry, err := h.tokenMaker.CreateRefreshToken(foundUser)
	if err != nil {
		slog.Error("failed to generate refresh token", "email", foundUser.Email, "error", err)
		return nil, "", "", ErrTokenCreationFailed
	}

	// hash the refresh token before storing it in the database
	hashedRefreshToken := utils.HashToken(refreshToken)

	newUserSession := token.Session{
		ID:           uuid.New(),
		UserEmail:    foundUser.Email,
		RefreshToken: hashedRefreshToken,
		IsRevoked:    false,
		CreatedAt:    refreshTokenIssuedAt,
		ExpiresAt:    refreshTokenExpiry,
	}

	err = h.sessionStore.CreateRefreshToken(ctx, &newUserSession)
	if err != nil {
		slog.Error("failed to create session", "email", foundUser.Email, "error", err)
		return nil, "", "", ErrSessionCreationFailed
	}

	return foundUser, accessToken, refreshToken, nil
}

func (h *authService) RefreshAccessToken(ctx context.Context, req RefreshTokenRequest) (string, error) {
	if req.RefreshToken == "" {
		return "", ErrInvalidRefreshToken
	}

	// hash the provided refresh token before comparing it with the stored hash
	hashedRefreshToken := utils.HashToken(req.RefreshToken)

	// retrieve the session from the database using the hashed refresh token
	retrievedSess, err := h.sessionStore.GetRefreshToken(ctx, hashedRefreshToken)
	if err != nil {
		if errors.Is(err, token.ErrNotFound) {
			return "", ErrInvalidRefreshToken
		}
		slog.Error("failed to retrieve session by refresh token", "refresh_token", req.RefreshToken, "error", err)
		return "", err
	}

	// check if the session is revoked or expired
	if retrievedSess.IsRevoked || time.Now().After(retrievedSess.ExpiresAt) {
		return "", ErrInvalidRefreshToken
	}

	// retrieve the user associated with the session
	foundUser, err := h.userStore.GetByEmail(ctx, retrievedSess.UserEmail)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			return "", ErrInvalidRefreshToken
		}
		slog.Error("failed to retrieve user by email", "email", retrievedSess.UserEmail, "error", err)
		return "", err
	}

	// generate a new access token for the user
	accessToken, err := h.tokenMaker.GenerateToken(foundUser)
	if err != nil {
		slog.Error("failed to generate access token", "email", foundUser.Email, "error", err)
		return "", err
	}

	return accessToken, nil
}

func (h *authService) LogoutUser(ctx context.Context, userId string) error {
	// check if the user exists
	foundUser, err := h.userStore.GetByID(ctx, userId)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			return ErrUserNotFound
		}
		slog.Error("failed to retrieve user by id", "id", userId, "error", err)
		return err
	}
	// check if the user has an active session and delete it
	err = h.sessionStore.DeleteSessionsByEmail(ctx, foundUser.Email)
	if err != nil {
		slog.Error("failed to delete sessions", "email", foundUser.Email, "error", err)
		return ErrFailedToDeleteSession
	}

	return nil
}
