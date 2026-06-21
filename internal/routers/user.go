package routers

import (
	"github.com/Milua25/go-job-application-tracker/internal/user"
	"github.com/gin-gonic/gin"
)

func registerUserRoutes(router *gin.RouterGroup, h *user.UserHandler, authMiddleware gin.HandlerFunc) {
	userGroup := router.Group("/", authMiddleware)
	userGroup.GET("/users", h.GetAllUsers)
	userGroup.GET("/users/:id", h.GetUserByID)
	// userGroup.POST("/users", h.CreateNewUser)
	userGroup.DELETE("/users/:id", h.DeleteUser)
}
