package auth

import (
	"errors"
	"log/slog"

	"github.com/Milua25/go-job-application-tracker/internal/authctx"
	"github.com/Milua25/go-job-application-tracker/internal/render"
	"github.com/Milua25/go-job-application-tracker/internal/token"
	"github.com/Milua25/go-job-application-tracker/internal/user"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	authService *authService
}

func actorFromContext(c *gin.Context) (Actor, error) {
	actor := Actor{
		UserID:    c.GetString(authctx.ContextKeyUID),
		Email:     c.GetString(authctx.ContextKeyEmail),
		SessionID: c.GetString(authctx.ContextKeySessionID),
		IsAdmin:   c.GetBool(authctx.ContextKeyIsAdmin),
	}

	if actor.UserID == "" || actor.Email == "" || actor.SessionID == "" {
		return Actor{}, errors.New("missing authenticated actor in context")
	}

	return actor, nil
}

// NewAuthHandler creates a new AuthHandler with the provided user store, session store, and token maker.
func NewAuthHandler(userStore user.Repository, sessionStore token.Repository, tokenMaker *token.JWTMaker) *AuthHandler {
	return &AuthHandler{
		authService: newAuthService(userStore, sessionStore, tokenMaker),
	}
}

// RegisterRoutes registers the authentication routes with the provided Gin router and authentication middleware.
func (h *AuthHandler) RegisterRoutes(r gin.IRouter, authMiddleware gin.HandlerFunc) {
	g := r.Group("/auth")
	g.POST("/register", h.RegisterUser)
	g.POST("/login", h.LoginUser)
	g.POST("/refresh", h.RefreshToken)

	// Protected routes that require authentication
	protected := g.Group("")
	protected.Use(authMiddleware)
	protected.POST("/logout", h.LogoutUser)
	protected.POST("/logout/all", h.DeleteAllSessionsForUser)
}

// RegisterUser handles user registration requests. It validates the request body, creates a new user, and returns the created user's details.
func (h *AuthHandler) RegisterUser(c *gin.Context) {
	// validate the request body
	var newUser RegisterRequest
	if err := c.ShouldBindJSON(&newUser); err != nil {
		render.ValidationError(c, err)
		return
	}
	createdUser, err := h.authService.register(c.Request.Context(), newUser)
	if err != nil {
		if errors.Is(err, ErrUserAlreadyExists) {
			render.ConflictResponseError(c, "email already registered", err)
			return
		}
		render.InternalServerError(c, "failed to register user", err)
		return
	}
	// Create the response object
	response := RegisterResponse{
		ID:        createdUser.ID.String(),
		FirstName: createdUser.FirstName,
		LastName:  createdUser.LastName,
		Email:     createdUser.Email,
		IsActive:  createdUser.IsActive,
		CreatedAt: createdUser.CreatedAt,
		UpdatedAt: createdUser.UpdatedAt,
	}

	render.Created(c, gin.H{
		"message": "user created successfully",
		"user":    response,
	})
}

// LoginUser handles user login requests. It validates the request body, authenticates the user, and returns the user's details along with access and refresh tokens.
func (h *AuthHandler) LoginUser(c *gin.Context) {
	// validate the request body
	var loginReq LoginRequest
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		render.ValidationError(c, err)
		return
	}

	foundUser, accessToken, accessTokenExpiresAt, refreshToken, refreshTokenExpiresAt, err := h.authService.loginUser(c.Request.Context(), loginReq)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			render.UnauthorizedBasicResponseError(c, "invalid email or password", err)
			return
		}
		render.InternalServerError(c, "failed to login user", err)
		return
	}

	// Create the response object
	response := LoginResponse{
		ID:                    foundUser.ID.String(),
		FirstName:             foundUser.FirstName,
		LastName:              foundUser.LastName,
		Email:                 foundUser.Email,
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  accessTokenExpiresAt,
		RefreshTokenExpiresAt: refreshTokenExpiresAt,
	}

	render.OK(c, gin.H{
		"message": "user logged in successfully",
		"user":    response,
	})
}

// LogoutUser handles user logout requests. It retrieves identity from context and logs out the current session.
func (h *AuthHandler) LogoutUser(c *gin.Context) {
	actor, err := actorFromContext(c)
	if err != nil {
		slog.Warn("authenticated actor not found in context")
		render.UnauthorizedBasicResponseError(c, "user not authenticated", err)
		return
	}
	slog.Info("active user found", "user_id", actor.UserID)

	if err = h.authService.LogoutUser(c.Request.Context(), actor.SessionID); err != nil {
		render.InternalServerError(c, "failed to logout user", err)
		return
	}

	render.OK(c, gin.H{
		"message": "user logged out successfully",
	})
}

// RefreshToken handles requests to refresh an access token. It validates the request body, generates a new access token using the provided refresh token, and returns the new access token along with its expiration time.
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var refreshReq RefreshTokenRequest
	if err := c.ShouldBindJSON(&refreshReq); err != nil {
		render.ValidationError(c, err)
		return
	}

	newAccessToken, accessTokenExpiresAt, err := h.authService.RefreshAccessToken(c.Request.Context(), refreshReq)
	if err != nil {
		if errors.Is(err, ErrInvalidRefreshToken) {
			render.UnauthorizedBasicResponseError(c, "invalid or expired refresh token", err)
			return
		}
		render.InternalServerError(c, "failed to refresh token", err)
		return
	}

	render.OK(c, gin.H{
		"message": "token refreshed successfully",
		"data": RenewAccessTokenResp{
			AccessToken:          newAccessToken,
			AccessTokenExpiresAt: accessTokenExpiresAt,
		},
	})
}

func (h *AuthHandler) DeleteAllSessionsForUser(c *gin.Context) {
	actor, err := actorFromContext(c)
	if err != nil {
		slog.Warn("authenticated actor not found in context")
		render.UnauthorizedBasicResponseError(c, "user not authenticated", err)
		return
	}

	userID, err := uuid.Parse(actor.UserID)
	if err != nil {
		render.UnauthorizedBasicResponseError(c, "invalid user id in token", err)
		return
	}

	if err = h.authService.DeleteAllSessionsForUser(c.Request.Context(), userID); err != nil {
		render.InternalServerError(c, "failed to delete all sessions for user", err)
		return
	}

	render.NoContent(c)
}
