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
	Status int `json:"-"`
	// Code    string `json:"code"`
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

// JSONError sends an error response.
func JSONError(c *gin.Context, status int, message string) {
	c.JSON(status, Response{
		Success: false,
		Error:   &ErrorInfo{Status: status, Message: message},
	})
}
