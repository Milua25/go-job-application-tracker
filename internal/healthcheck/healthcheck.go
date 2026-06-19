package healthcheck

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type HealthCheckResponse struct {
	Status string `json:"status"`
	Checks gin.H  `json:"checks"`
}

// HealthCheckHandler handles health check requests.
type HealthCheckHandler struct {
	// DB *sql.DB
}

func (*HealthCheckHandler) CheckHealth(c *gin.Context) {

	_, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	db_healthy := true
	redis_healthy := true

	// Example database check (uncomment and implement as needed)
	// if err := h.DB.PingContext(ctx); err != nil {
	// 	checks["database"] = "unhealthy"
	// 	healthy = false
	// } else {
	// 	checks["database"] = "healthy"
	// }
	if db_healthy && redis_healthy {
		health := HealthCheckResponse{
			Status: "healthy",
			Checks: gin.H{
				"database": "up",
				"redis":    "up",
			},
		}
		// Send a 200 OK response with the health status
		c.JSON(http.StatusOK, health)
	}

}
