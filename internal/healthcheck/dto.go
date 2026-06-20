package healthcheck

import "github.com/gin-gonic/gin"

type HealthCheckResponse struct {
	Status string `json:"status"`
	Checks gin.H  `json:"checks"`
}
