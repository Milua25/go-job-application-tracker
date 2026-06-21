package routers

import (
	"github.com/Milua25/go-job-application-tracker/internal/auth"
	"github.com/gin-gonic/gin"
)

func registerAuthRoutes(router *gin.RouterGroup, h *auth.AuthHandler) {
	authGroup := router.Group("/")
	authGroup.POST("/auth/register", h.RegisterUser)
}
