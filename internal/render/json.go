package render

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
	// Meta    *Meta       `json:"meta,omitempty"`
}

type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

//	type Meta struct {
//		Page       int `json:"page,omitempty"`
//		PerPage    int `json:"per_page,omitempty"`
//		Total      int `json:"total,omitempty"`
//		TotalPages int `json:"total_pages,omitempty"`
//	}
//

// OK sends a success response.
func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
}

func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// JSONError sends an error response.
func JSONError(c *gin.Context, status int, code, message string) {
	c.JSON(status, Response{
		Success: false,
		Error:   &ErrorInfo{Code: code, Message: message},
	})
}

// Fail sends a non-2xx response that still carries structured data (e.g. health checks).
func Fail(c *gin.Context, status int, data interface{}) {
	c.JSON(status, Response{
		Success: false,
		Data:    data,
	})
}
