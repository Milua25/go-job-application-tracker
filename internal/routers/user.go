package routers

import (
	"github.com/Milua25/go-job-application-tracker/internal/user"
	"github.com/gin-gonic/gin"
)

func registerUserRoutes(router *gin.RouterGroup) {

	// router.GET("/ping", func(c *gin.Context) {
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"message": "pong",
	// 	})
	// })

	// User routes
	userGroup := router.Group("/")
	userGroup.GET("/users", user.GetAllUsers)
	userGroup.GET("/users/:id", user.GetUserByID)
}
