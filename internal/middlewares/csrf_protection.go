package middlewares

import (
	"net/http"

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
	return sessions.Sessions("mysession", store)
}
