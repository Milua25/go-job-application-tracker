package middlewares

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Cors() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:              []string{"http://localhost:8000"},
		AllowMethods:              []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:              []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:             []string{"Content-Length"},
		AllowCredentials:          true,
		MaxAge:                    12 * time.Hour,
		OptionsResponseStatusCode: http.StatusNoContent,
	})
}

func RequireOrigin() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodOptions && c.GetHeader("Origin") == "" {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		c.Next()
	}
}
