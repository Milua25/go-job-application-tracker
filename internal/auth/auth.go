package auth

import (
	"errors"
	"time"

	"github.com/Milua25/go-job-application-tracker/internal/render"
	"github.com/Milua25/go-job-application-tracker/internal/user"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	userStore   Repository
	authService *AuthService
}

func NewAuthHandler(store Repository, secretKey, issuer, expireIn, refreshExpireIn string) *AuthHandler {
	return &AuthHandler{
		userStore:   store,
		authService: newAuthService(secretKey, issuer, expireIn, refreshExpireIn),
	}
}

func (h *AuthHandler) Service() *AuthService {
	return h.authService
}

func (h *AuthHandler) RegisterUser(c *gin.Context) {
	// parse request body into User struct
	var newUser RegisterRequest
	if err := c.ShouldBindJSON(&newUser); err != nil {
		render.ValidationError(c, err)
		return
	}
	// check if user with the same email already exists
	existingUser, err := h.userStore.GetByEmail(c.Request.Context(), newUser.Email)
	if err != nil {
		render.InternalServerError(c, "failed to check existing user", err)
		return
	}
	if existingUser != nil {
		render.ConflictResponseError(c, "user with this email already exists", errors.New("email already registered"))
		return
	}

	// create new user and add to users slice
	createdUser := user.User{
		ID:           uuid.New(),
		Email:        newUser.Email,
		PasswordHash: "hashedpassword", // In real implementation, hash the password properly
		FirstName:    newUser.FirstName,
		LastName:     newUser.LastName,
		Timezone:     "UTC", // Default timezone, can be updated later
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// create new user in database
	err = h.userStore.Create(c.Request.Context(), &createdUser)
	if err != nil {
		render.InternalServerError(c, "failed to create user", err)
		return
	}

	//return user created response
	response := RegisterResponse{
		ID:        createdUser.ID.String(),
		FirstName: createdUser.FirstName,
		LastName:  createdUser.LastName,
		Email:     createdUser.Email,
		IsActive:  createdUser.IsActive,
		CreatedAt: createdUser.CreatedAt,
		UpdatedAt: createdUser.UpdatedAt,
	}

	render.OK(c, gin.H{
		"message": "user created successfully",
		"user":    response,
	})
}

func (h *AuthHandler) LoginUser(c *gin.Context) {
	// parse request body into LoginRequest struct
	var loginReq LoginRequest
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		render.ValidationError(c, err)
		return
	}
	// find user by email
	user, err := h.userStore.GetByEmail(c.Request.Context(), loginReq.Email)
	if err != nil {
		render.InternalServerError(c, "failed to retrieve user", err)
		return
	}
	if user == nil || user.PasswordHash != "hashedpassword" { // In real implementation, compare hashed passwords properly
		render.UnAuthorizedBasicResponseError(c, "invalid email or password", errors.New("authentication failed"))
		return
	}

	// token generation logic would go here (e.g., JWT token creation)
	token, refreshToken, err := h.authService.GenerateToken(user)
	if err != nil {
		render.InternalServerError(c, "failed to generate token", err)
		return
	}
	// return login response
	response := LoginResponse{
		ID:           user.ID.String(),
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Email:        user.Email,
		Token:        token,
		RefreshToken: refreshToken,
	}

	render.OK(c, gin.H{
		"message": "login successful",
		"user":    response,
	})
}

func (h *AuthHandler) CreateRefreshToken(c *gin.Context) {
	// In a real implementation, you would validate the refresh token, check if it's still valid, and then generate a new access token (and possibly a new refresh token)
	var refreshReq RefreshRequest
	if err := c.ShouldBindJSON(&refreshReq); err != nil {
		render.ValidationError(c, err)
		return
	}
	// validate refresh token and generate new access token logic would go here

	// For demonstration, we'll just return a new token without actual validation
}

func (h *AuthHandler) LogoutUser(c *gin.Context) {

}
