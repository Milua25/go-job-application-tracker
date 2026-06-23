package auth

import (
	"errors"

	"github.com/Milua25/go-job-application-tracker/internal/render"
	"github.com/Milua25/go-job-application-tracker/internal/session"
	"github.com/Milua25/go-job-application-tracker/internal/token"
	"github.com/Milua25/go-job-application-tracker/internal/user"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *authService
}

func NewAuthHandler(userStore user.Repository, sessionStore session.Repository, tokenMaker *token.JWTMaker) *AuthHandler {
	return &AuthHandler{
		authService: newAuthService(userStore, sessionStore, tokenMaker),
	}
}

func (h *AuthHandler) RegisterRoutes(r gin.IRouter) {
	g := r.Group("/auth")
	g.POST("/register", h.RegisterUser)
	g.POST("/login", h.LoginUser)
}

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

func (h *AuthHandler) LoginUser(c *gin.Context) {
	var loginReq LoginRequest
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		render.ValidationError(c, err)
		return
	}

	foundUser, accessToken, refreshToken, err := h.authService.loginUser(c.Request.Context(), loginReq)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			render.UnauthorizedBasicResponseError(c, "invalid email or password", err)
			return
		}
		render.InternalServerError(c, "failed to login user", err)
		return
	}

	// TODO: Delete old sessions for the user before creating new ones to prevent unique constraint violations on token
	// This should be implemented in authService.loginUser() - add call to sessionStore.DeleteSessionsByUserEmail()

	// Create the response object
	response := LoginResponse{
		ID:                    foundUser.ID.String(),
		FirstName:             foundUser.FirstName,
		LastName:              foundUser.LastName,
		Email:                 foundUser.Email,
		Token:                 accessToken,
		RefreshToken:          refreshToken,
		TokenExpiresAt:        foundUser.CreatedAt,
		RefreshTokenExpiresAt: foundUser.UpdatedAt,
	}

	render.OK(c, gin.H{
		"message": "user logged in successfully",
		"user":    response,
	})
}

// func (h *AuthHandler) RefreshToken(c *gin.Context) {
// 	var refreshReq RefreshRequest
// 	if err := c.ShouldBindJSON(&refreshReq); err != nil {
// 		render.ValidationError(c, err)
// 		return
// 	}

// 	// Call the authService to handle the refresh token logic
// 	newAccessToken, err := h.authService.createRefreshToken(c.Request.Context(), refreshReq)
// 	if err != nil {
// 		render.InternalServerError(c, "failed to refresh token", err)
// 		return
// 	}

// 	render.OK(c, gin.H{
// 		"message":      "token refreshed successfully",
// 		"access_token": newAccessToken,
// 	})
// }
