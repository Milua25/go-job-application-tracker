package middlewares

import (
	"compress/gzip"

	zip "github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func Compression() gin.HandlerFunc {

	zipHandlerfunc := zip.Gzip(gzip.DefaultCompression, zip.WithExcludedPaths([]string{"/api/v1/health"}))

	return zipHandlerfunc
	// return func(c *gin.Context) {
	// 	// check if the client supports gzip encoding
	// 	// fix
	// 	if !strings.Contains(c.Request.Header.Get("Accept-Encoding"), "gzip") {
	// 		c.Next()
	// 		return
	// 	}

	// 	gz := gzip.NewWriter(c.Writer)
	// 	defer gz.Close()

	// 	w := &gzipResponseWriter{ResponseWriter: c.Writer, Writer: gz}
	// 	c.Writer = w
	// 	c.Header("Content-Encoding", "gzip")

	// 	c.Next()
	// }
}

type gzipResponseWriter struct {
	gin.ResponseWriter
	Writer *gzip.Writer
}

func (w *gzipResponseWriter) Write(data []byte) (int, error) {
	return w.Writer.Write(data)
}
