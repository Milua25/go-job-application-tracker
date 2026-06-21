package middleware

import (
	"log"
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
)

type HPPOptions struct {
	CheckBody        bool
	CheckQuery       bool
	CheckContentType string
	Whitelist        []string
}

func (opts HPPOptions) Hpp() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		if opts.CheckBody && (method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch || method == http.MethodDelete) && isCorrectContentType(c, opts.CheckContentType) {
			// filter body parameters
			filterParameters(c, opts.Whitelist)
		}
		if opts.CheckQuery && c.Request.URL.Query() != nil {
			// filter query parameters
			filterQueryParameters(c, opts.Whitelist)
		}
		c.Next()
	}
}

func isCorrectContentType(c *gin.Context, contentType string) bool {
	return strings.Contains(c.Request.Header.Get("Content-Type"), contentType)
}

func filterParameters(c *gin.Context, whitelist []string) {
	err := c.Request.ParseForm()
	if err != nil {
		log.Printf("Error parsing form: %v", err)
		return
	}

	for key, values := range c.Request.PostForm {
		if len(values) > 1 {
			c.Request.PostForm.Set(key, values[0])
		}
		if !slices.Contains(whitelist, key) {
			c.Request.PostForm[key] = []string{values[0]}
		}
	}
}

// filterQueryParameters filters query parameters based on the whitelist
func filterQueryParameters(c *gin.Context, whitelist []string) {
	query := c.Request.URL.Query()
	for key, values := range query {
		if len(values) > 1 {
			query.Set(key, values[0])
		}
		if !slices.Contains(whitelist, key) {
			delete(query, key)
		}
	}
	c.Request.URL.RawQuery = query.Encode()
}
