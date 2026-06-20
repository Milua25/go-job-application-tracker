package healthcheck

import (
	"context"
	"net/http"
	"time"

	"github.com/Milua25/go-job-application-tracker/internal/render"
	"github.com/gin-gonic/gin"
)


// HealthCheckHandler handles health check requests.
type HealthCheckHandler struct {
	// DB *sql.DB
}

func (*HealthCheckHandler) CheckHealth(c *gin.Context) {

	_, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	db_healthy := false
	redis_healthy := true
	dbStatus := "down"
	redisStatus := "up"

	if !db_healthy {
		dbStatus = "down"
	}

	if !redis_healthy {
		redisStatus = "down"
	}

	// Example database check (uncomment and implement as needed)
	// if err := h.DB.PingContext(ctx); err != nil {
	// 	checks["database"] = "unhealthy"
	// 	healthy = false
	// } else {
	// 	checks["database"] = "healthy"
	// }
	health := HealthCheckResponse{
		Checks: gin.H{
			"database": dbStatus,
			"redis":    redisStatus,
		},
	}

	if db_healthy && redis_healthy {
		health.Status = "healthy"
		render.OK(c, health)
		return
	}

	health.Status = "unhealthy"
	render.Fail(c, http.StatusServiceUnavailable, health)
}
