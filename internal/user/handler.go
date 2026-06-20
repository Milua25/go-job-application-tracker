package user

import (
	"github.com/Milua25/go-job-application-tracker/internal/render"
	"github.com/gin-gonic/gin"
)

func GetUser(c *gin.Context) {

	render.OK(c, gin.H{
		"message": "get user",
		// "csrf":    csrfToken,
	})
}
