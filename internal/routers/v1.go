package routers

import "github.com/gin-gonic/gin"

func APIv1Routes(router *gin.Engine) {
	v1 := router.Group("/api/v1")

	registerUserRoutes(v1)
	registerHealthCheckRoutes(v1)
}
