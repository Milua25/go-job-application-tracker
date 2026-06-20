package healthcheck

import (
	"github.com/Milua25/go-job-application-tracker/internal/render"
	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	render.OK(c, gin.H{"message": "pong"})
}
