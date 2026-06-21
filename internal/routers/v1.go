package routers

import (
	"github.com/Milua25/go-job-application-tracker/internal/auth"
	"github.com/Milua25/go-job-application-tracker/internal/healthcheck"
	"github.com/Milua25/go-job-application-tracker/internal/user"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	User           *user.UserHandler
	HealthCheck    *healthcheck.HealthCheckHandler
	Auth           *auth.AuthHandler
	AuthMiddleware gin.HandlerFunc
}

func APIv1Routes(router *gin.Engine, h Handlers) {
	v1 := router.Group("/api/v1")

	registerUserRoutes(v1, h.User, h.AuthMiddleware)
	registerHealthCheckRoutes(v1, h.HealthCheck)
	registerAuthRoutes(v1, h.Auth)
}
