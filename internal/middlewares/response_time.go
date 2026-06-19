package middlewares

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func ResponseTime() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Record the start time
		start := time.Now()

		wrappedWriter := &responseTimeWriter{ResponseWriter: c.Writer, statusCode: http.StatusOK}
		c.Writer = wrappedWriter
		c.Next()

		elapsed := time.Since(start)
		ms := strconv.FormatFloat(elapsed.Seconds()*1000, 'f', 3, 64)
		log.Printf("Method: %s, Request: %s %s, Response Time: %sms, Status: %d", c.Request.Method, c.Request.URL.EscapedPath(), c.Request.Method, ms, wrappedWriter.statusCode)
	}
}

type responseTimeWriter struct {
	gin.ResponseWriter
	statusCode int
}

func (w *responseTimeWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}
