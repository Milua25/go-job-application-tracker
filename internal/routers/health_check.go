package routers

import (
	"github.com/Milua25/go-job-application-tracker/internal/healthcheck"
	"github.com/gin-gonic/gin"
)

func registerHealthCheckRoutes(router *gin.RouterGroup, h *healthcheck.HealthCheckHandler) {
	healthGroup := router.Group("/")
	healthGroup.GET("/health", h.CheckHealth)
	healthGroup.GET("/ping", h.Ping)
}
