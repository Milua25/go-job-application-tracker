package render

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Document is the JSON:API top-level document.
// Exactly one of Data or Meta must be set per the spec.
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

// Fail sends a non-2xx response that still carries structured data (e.g. health checks).
func Fail(c *gin.Context, status int, data interface{}) {
	c.JSON(status, Document{Data: data})
}
