package routers

import (
	"github.com/Milua25/go-job-application-tracker/internal/auth"
	"github.com/Milua25/go-job-application-tracker/internal/healthcheck"
	"github.com/Milua25/go-job-application-tracker/internal/user"
	"github.com/gin-gonic/gin"
)

func RegisterV1Routes(router *gin.Engine, authMiddleware gin.HandlerFunc, u *user.UserHandler, hc *healthcheck.HealthCheckHandler, a *auth.AuthHandler) {
	v1 := router.Group("/api/v1")

	u.RegisterRoutes(v1, authMiddleware)
	hc.RegisterRoutes(v1)
	a.RegisterRoutes(v1)
}
