package middlewares

import (
	"net/http"
	"time"

	"github.com/Milua25/go-job-application-tracker/internal/render"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

func CSRFProtection() gin.HandlerFunc {
	return csrf.Middleware(csrf.Options{
		Secret: "secret123",
		ErrorFunc: func(ctx *gin.Context) {
			render.JSONError(ctx, http.StatusBadRequest, "CSRF token mismatch")
		},
	})
}

func SessionCSRFProtection() gin.HandlerFunc {
	store := cookie.NewStore([]byte("secret"))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   int((24 * time.Hour).Seconds()),
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
	})
	return sessions.Sessions("mysession", store)
}

// CSRFTokenHeader generates the CSRF token on GET requests and exposes it
// via the X-CSRF-TOKEN response header so clients can use it in subsequent
// mutating requests.
func CSRFTokenHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodGet {
			c.Header("X-CSRF-TOKEN", csrf.GetToken(c))
		}
		c.Next()
	}
}
