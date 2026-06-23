package render

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Document is the JSON:API top-level document.
// Per JSON:API spec, exactly one of Data or Meta should be set at the top level.
// The struct allows both to be omitted or both to be set; callers must enforce this constraint.
type Document struct {
	Data  interface{} `json:"data,omitempty"`
	Meta  interface{} `json:"meta,omitempty"`
	Links interface{} `json:"links,omitempty"`
}

// OK sends a 200 response with primary data.
func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Document{Data: data})
}

// OKMeta sends a 200 response with only a meta object (no primary data).
func OKMeta(c *gin.Context, meta interface{}) {
	c.JSON(http.StatusOK, Document{Meta: meta})
}

// NoContent sends a 204 with no body.
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Document{Data: data})
}

// Fail sends a non-2xx response that still carries structured data (e.g. health checks).
// status must be a valid HTTP status code (100-599); behavior is undefined for invalid codes.
func Fail(c *gin.Context, status int, data interface{}) {
	if status < 100 || status > 599 {
		status = http.StatusInternalServerError
	}
	c.JSON(status, Document{Data: data})
}
