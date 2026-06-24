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
	if err != nil && !errors.Is(err, user.ErrNotFound) {
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

func (h *authService) loginUser(ctx context.Context, req LoginRequest) (*user.User, string, time.Time, string, time.Time, error) {
	slog.Debug("logging in user", "email", req.Email)

	foundUser, err := h.userStore.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			return nil, "", time.Time{}, "", time.Time{}, ErrInvalidCredentials
		}
		slog.Error("failed to retrieve user by email", "email", req.Email, "error", err)
		return nil, "", time.Time{}, "", time.Time{}, err
	}

	if !utils.VerifyPassword(foundUser.PasswordHash, req.Password) {
		slog.Warn("authentication failed", "email", req.Email)
		return nil, "", time.Time{}, "", time.Time{}, ErrInvalidCredentials
	}

	sessionID := uuid.New()

	accessToken, accessTokenExpiresAt, err := h.tokenMaker.GenerateToken(foundUser, sessionID.String())
	if err != nil {
		slog.Error("failed to generate token", "email", foundUser.Email, "error", err)
		return nil, "", time.Time{}, "", time.Time{}, err
	}

	refreshToken, refreshTokenIssuedAt, refreshTokenExpiresAt, err := h.tokenMaker.CreateRefreshToken(foundUser)
	if err != nil {
		slog.Error("failed to generate refresh token", "email", foundUser.Email, "error", err)
		return nil, "", time.Time{}, "", time.Time{}, ErrTokenCreationFailed
	}

	newUserSession := token.Session{
		ID:           sessionID,
		UserID:       foundUser.ID,
		RefreshToken: utils.HashToken(refreshToken),
		IsRevoked:    false,
		CreatedAt:    refreshTokenIssuedAt,
		ExpiresAt:    refreshTokenExpiresAt,
	}

	if err = h.sessionStore.CreateSession(ctx, &newUserSession); err != nil {
		slog.Error("failed to create session", "email", foundUser.Email, "error", err)
		return nil, "", time.Time{}, "", time.Time{}, ErrSessionCreationFailed
	}

	return foundUser, accessToken, accessTokenExpiresAt, refreshToken, refreshTokenExpiresAt, nil
}

func (h *authService) RefreshAccessToken(ctx context.Context, req RefreshTokenRequest) (string, time.Time, error) {
	if req.RefreshToken == "" {
		return "", time.Time{}, ErrInvalidRefreshToken
	}

	retrievedSess, err := h.sessionStore.GetSessionByRefreshToken(ctx, utils.HashToken(req.RefreshToken))
	if err != nil {
		if errors.Is(err, token.ErrNotFound) {
			return "", time.Time{}, ErrInvalidRefreshToken
		}
		slog.Error("failed to retrieve session by refresh token", "error", err)
		return "", time.Time{}, err
	}

	if retrievedSess.IsRevoked || time.Now().After(retrievedSess.ExpiresAt) {
		return "", time.Time{}, ErrInvalidRefreshToken
	}

	foundUser, err := h.userStore.GetByID(ctx, retrievedSess.UserID.String())
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			return "", time.Time{}, ErrInvalidRefreshToken
		}
		slog.Error("failed to retrieve user by id", "user_id", retrievedSess.UserID, "error", err)
		return "", time.Time{}, err
	}

	accessToken, accessTokenExpiresAt, err := h.tokenMaker.GenerateToken(foundUser, retrievedSess.ID.String())
	if err != nil {
		slog.Error("failed to generate access token", "email", foundUser.Email, "error", err)
		return "", time.Time{}, err
	}

	return accessToken, accessTokenExpiresAt, nil
}

func (h *authService) LogoutUser(ctx context.Context, sessionID string) error {
	if err := h.sessionStore.DeleteSessionByID(ctx, sessionID); err != nil {
		if errors.Is(err, token.ErrNotFound) {
			return ErrNotFound
		}
		slog.Error("failed to delete session by id", "id", sessionID, "error", err)
		return err
	}
	return nil
}

func (h *authService) DeleteAllSessionsForUser(ctx context.Context, userID uuid.UUID) error {
	if err := h.sessionStore.DeleteSessionsByUserID(ctx, userID); err != nil {
		slog.Error("failed to delete sessions for user", "user_id", userID, "error", err)
		return err
	}
	return nil
}
