package auth

import (
	"errors"
	"time"

	"github.com/Milua25/go-job-application-tracker/internal/render"
	"github.com/Milua25/go-job-application-tracker/internal/user"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TokenGenerator interface {
	GenerateToken(firstName, lastName, email, uid string) (string, string, error)
}

type AuthHandler struct {
	store Repository
	token TokenGenerator
}

func NewAuthHandler(store Repository, token TokenGenerator) *AuthHandler {
	return &AuthHandler{store: store, token: token}
}

func (h *AuthHandler) RegisterUser(c *gin.Context) {
	// parse request body into User struct
	var newUser RegisterRequest
	if err := c.ShouldBindJSON(&newUser); err != nil {
		render.BadRequestError(c, "invalid request body", err)
		return
	}
	// check if user with the same email already exists
	existingUser, err := h.store.GetByEmail(c.Request.Context(), newUser.Email)
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
	err = h.store.Create(c.Request.Context(), &createdUser)
	if err != nil {
		render.InternalServerError(c, "failed to create user", err)
		return
	}

	// token generation logic would go here (e.g., JWT token creation)
	token, refreshToken, err := h.token.GenerateToken(createdUser.FirstName, createdUser.LastName, createdUser.Email, createdUser.ID.String())
	if err != nil {
		render.InternalServerError(c, "failed to generate token", err)
		return
	}
	//return user created response
	response := RegisterResponse{
		ID:           createdUser.ID.String(),
		FirstName:    createdUser.FirstName,
		LastName:     createdUser.LastName,
		Email:        createdUser.Email,
		Token:        token,
		RefreshToken: refreshToken,
		IsActive:     createdUser.IsActive,
		CreatedAt:    createdUser.CreatedAt,
		UpdatedAt:    createdUser.UpdatedAt,
	}

	render.OK(c, gin.H{
		"message": "user created successfully",
		"user":    response,
	})
}
