package routers

import (
	"github.com/Milua25/go-job-application-tracker/internal/healthcheck"
	"github.com/gin-gonic/gin"
)

func registerHealthCheckRoutes(router *gin.RouterGroup) {
	healthGroup := router.Group("/")
	healthGroup.GET("/health", (&healthcheck.HealthCheckHandler{}).CheckHealth)
}
