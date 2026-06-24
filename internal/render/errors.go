package render

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/Milua25/go-job-application-tracker/pkg/utils"
	"github.com/gin-gonic/gin"
)

type errorObject struct {
	Status string `json:"status"`
	Code   string `json:"code"`
	Detail string `json:"detail"`
}

type errorsDocument struct {
	Errors []errorObject `json:"errors"`
}

// JSONError sends an error response following the JSON:API errors document structure
// and aborts the handler chain.
func JSONError(c *gin.Context, status int, code, detail string) {
	c.AbortWithStatusJSON(status, errorsDocument{
		Errors: []errorObject{{
			Status: http.StatusText(status),
			Code:   code,
			Detail: detail,
		}},
	})
}

// ValidationError extracts field-level validation errors and sends a 422 response
// with one error object per failing field.
func ValidationError(c *gin.Context, err error) {
	slog.Debug("validation error", "error", err)
	fieldErrors := utils.ExtractValidationErrors(err)
	var errs []errorObject
	if len(fieldErrors) == 0 {
		// If no field errors were extracted, send a generic validation error response
		errs = []errorObject{{
			Status: http.StatusText(http.StatusUnprocessableEntity),
			Code:   "VALIDATION_ERROR",
			Detail: err.Error(),
		}}
	} else {
		errs = make([]errorObject, len(fieldErrors))
		for i, fe := range fieldErrors {
			errs[i] = errorObject{
				Status: http.StatusText(http.StatusUnprocessableEntity),
				Code:   "VALIDATION_ERROR",
				Detail: fe.Message,
			}
		}
	}
	c.AbortWithStatusJSON(http.StatusUnprocessableEntity, errorsDocument{Errors: errs})
}

// InternalServerError logs and sends a 500 response.
// Logging uses stdlib log for simplicity; can be migrated to structured logging if needed.
func InternalServerError(c *gin.Context, msg string, err error) {
	slog.Error("internal server error", "method", c.Request.Method, "path", c.Request.URL.Path, "error", err)
	JSONError(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", msg)
}

// badRequestError logs and sends a 400 response.
func BadRequestError(c *gin.Context, msg string, err error) {
	slog.Error("bad request error", "method", c.Request.Method, "path", c.Request.URL.Path, "error", err)
	JSONError(c, http.StatusBadRequest, "BAD_REQUEST", msg)
}

// notFoundError logs and sends a 404 response.
func NotFoundError(c *gin.Context, msg string) {
	slog.Error("not found error", "method", c.Request.Method, "path", c.Request.URL.Path)
	JSONError(c, http.StatusNotFound, "NOT_FOUND", msg)
}

// conflictResponseError logs and sends a 409 response.
func ConflictResponseError(c *gin.Context, msg string, err error) {
	slog.Error("conflict server error", "method", c.Request.Method, "path", c.Request.URL.Path, "error", err)
	JSONError(c, http.StatusConflict, "CONFLICT", msg)
}

// UnauthorizedResponseError logs and sends a 401 response.
func UnauthorizedResponseError(c *gin.Context, msg string, err error) {
	slog.Error("unauthorized server error", "method", c.Request.Method, "path", c.Request.URL.Path, "error", err)
	JSONError(c, http.StatusUnauthorized, "UNAUTHORIZED", msg)
}

// UnauthorizedBasicResponseError logs and sends a 401 response.
func UnauthorizedBasicResponseError(c *gin.Context, msg string, err error) {
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
