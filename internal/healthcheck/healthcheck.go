package healthcheck

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/Milua25/go-job-application-tracker/internal/render"
	"github.com/gin-gonic/gin"
)

type HealthCheckHandler struct {
	db *sql.DB
}

// type HealthCheckHandler struct {
// 	db *sql.DB
// }

//	func NewHealthCheckHandler(db *sql.DB) *HealthCheckHandler {
//		return &HealthCheckHandler{db: db}
func NewHealthCheckHandler(db *sql.DB) *HealthCheckHandler {
	return &HealthCheckHandler{db: db}
}

func (h *HealthCheckHandler) CheckHealth(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	dbHealthy := h.db.PingContext(ctx) == nil

	dbStatus := "down"
	if dbHealthy {
		dbStatus = "up"
	}

	health := HealthCheckResponse{
		Checks: gin.H{"database": dbStatus},
	}

	if dbHealthy {
		health.Status = "healthy"
		render.OK(c, health)
		return
	}

	health.Status = "unhealthy"
	render.Fail(c, http.StatusServiceUnavailable, health)
}

func (h *HealthCheckHandler) Ping(c *gin.Context) {
	render.OK(c, gin.H{"message": "pong"})
}
