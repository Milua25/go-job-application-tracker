package render

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// internalServerError logs and sends a 500 response.
func InternalServerError(c *gin.Context, msg string, err error) {
	log.Printf("internal server error: %s path: %s error: %s", c.Request.Method, c.Request.URL.Path, err.Error())
	JSONError(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", msg)
}

// badRequestError logs and sends a 400 response.
func BadRequestError(c *gin.Context, msg string, err error) {
	log.Printf("bad request error: %s path: %s error: %s", c.Request.Method, c.Request.URL.Path, err.Error())
	JSONError(c, http.StatusBadRequest, "BAD_REQUEST", msg)
}

// notFoundError logs and sends a 404 response.
func NotFoundError(c *gin.Context, msg string) {
	log.Printf("not found error: %s path: %s", c.Request.Method, c.Request.URL.Path)
	JSONError(c, http.StatusNotFound, "NOT_FOUND", msg)
}

// conflictResponseError logs and sends a 409 response.
func ConflictResponseError(c *gin.Context, msg string, err error) {
	log.Printf("conflict server error: %s path: %s error: %s", c.Request.Method, c.Request.URL.Path, err.Error())
	JSONError(c, http.StatusConflict, "CONFLICT", msg)
}

// unAuthorizedResponseError logs and sends a 401 response.
func UnAuthorizedResponseError(c *gin.Context, msg string, err error) {
	log.Printf("unauthorized server error: %s path: %s error: %s", c.Request.Method, c.Request.URL.Path, err.Error())
	JSONError(c, http.StatusUnauthorized, "UNAUTHORIZED", msg)
}

// unAuthorizedBasicResponseError logs and sends a 401 response.
func UnAuthorizedBasicResponseError(c *gin.Context, msg string, err error) {
	log.Printf("unauthorized server error: %s path: %s error: %s", c.Request.Method, c.Request.URL.Path, err.Error())
	c.Header("www-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
	JSONError(c, http.StatusUnauthorized, "UNAUTHORIZED", msg)
}

// forbiddenResponseError logs and sends a 403 response.
func ForbiddenResponseError(c *gin.Context, msg string, err error) {
	log.Printf("forbidden response error: %s path: %s error: %s", c.Request.Method, c.Request.URL.Path, err.Error())
	JSONError(c, http.StatusForbidden, "FORBIDDEN", msg)
}

// rateLimitExceededResponse logs and sends a 429 response.
func RateLimitExceededResponse(c *gin.Context, retryAfter string) {
	log.Printf("rate limit exceeded: %s path: %s", c.Request.Method, c.Request.URL.Path)
	c.Header("Retry-After", retryAfter)
	JSONError(c, http.StatusTooManyRequests, "RATE_LIMIT_EXCEEDED", "rate limit exceeded, retry after: "+retryAfter)
}
