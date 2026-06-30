package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set("request_id", requestID)

		correlationID := c.GetHeader("X-Correlation-ID")
		if correlationID == "" {
			correlationID = requestID
		}
		c.Set("correlation_id", correlationID)

		c.Header("X-Request-ID", requestID)
		c.Header("X-Correlation-ID", correlationID)
		c.Next()
	}
}
