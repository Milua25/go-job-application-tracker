package healthcheck

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"time"

	"github.com/Milua25/go-job-application-tracker/internal/render"
	"github.com/gin-gonic/gin"
)

type HealthCheckHandler struct {
	db *sql.DB
}

func NewHealthCheckHandler(db *sql.DB) *HealthCheckHandler {
	return &HealthCheckHandler{db: db}
}

func (h *HealthCheckHandler) CheckHealth(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	dbStatus := "up"
	if err := h.db.PingContext(ctx); err != nil {
		slog.Error("database health check failed", "error", err)
		dbStatus = "down"
	}

	health := HealthCheckResponse{
		Status: func() string {
			if dbStatus == "up" {
				return "healthy"
			}
			return "unhealthy"
		}(),
		Checks: gin.H{"database": dbStatus},
	}

	slog.Info("health check result", "status", health.Status)

	if dbStatus == "down" {
		render.Fail(c, http.StatusServiceUnavailable, health)
		return
	}

	render.OK(c, health)
}

func (h *HealthCheckHandler) Ping(c *gin.Context) {
	render.OK(c, gin.H{"message": "pong"})
}
