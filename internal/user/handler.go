package user

import (
	"log"
	"time"

	"github.com/Milua25/go-job-application-tracker/internal/render"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var users = []User{
	{
		ID:           uuid.MustParse("6f78a05d-32b8-4b0e-8d5e-62e984a5969f"),
		Email:        "ayo@example.com",
		PasswordHash: "hashedpassword",
		FullName:     "Ayo Johnson",
		Timezone:     "America/New_York",
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	},
	{
		ID:           uuid.MustParse("1dd54345-3f2d-464e-98ab-1a7b982793dc"), // uuid.New(),
		Email:        "tayo@example.com",
		PasswordHash: "hashedpassword",
		FullName:     "Tayo Smith",
		Timezone:     "America/Los_Angeles",
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	},
}

func GetAllUsers(c *gin.Context) {
	list_of_users := users
	// database lookup for all users
	render.OK(c, gin.H{
		"message": "get all users",
		"users":   list_of_users,
	})
}

func GetUserByID(c *gin.Context) {
	id := c.Param("id")

	// check the id
	if _, err := uuid.Parse(id); err != nil {
		render.BadRequestError(c, "invalid user id", err)
		return
	}

	found_user := User{}

	// database lookup for user by id
	for _, user := range users {
		if user.ID.String() == id {
			log.Printf("found user: %s", user.Email)
			found_user = user
		}
	}
	// if user not found, return 404
	if found_user.ID == uuid.Nil {
		render.NotFoundError(c, "user not found")
		return
	}

	// return the found user
	render.OK(c, gin.H{
		"message": "get user by id",
		"user":    found_user,
	})
}
