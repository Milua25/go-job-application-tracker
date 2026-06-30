package healthcheck

import (
	"database/sql"
	"net/http"

	"github.com/Milua25/go-job-application-tracker/internal/render"
	"github.com/gin-gonic/gin"
)

type HealthCheckHandler struct {
	svc *healthCheckService
}

func NewHealthCheckHandler(db *sql.DB) *HealthCheckHandler {
	return &HealthCheckHandler{svc: newHealthCheckService(db)}
}

func (h *HealthCheckHandler) RegisterRoutes(r gin.IRouter) {
	g := r.Group("/")
	g.GET("/health", h.CheckHealth)
	g.GET("/ping", h.Ping)
}

func (h *HealthCheckHandler) CheckHealth(c *gin.Context) {
	status := h.svc.check(c.Request.Context())

	dbStatus := "up"
	if !status.dbUp {
		dbStatus = "down"
	}

	response := HealthCheckResponse{
		Status: func() string {
			if status.healthy {
				return "healthy"
			}
			return "unhealthy"
		}(),
		Checks: gin.H{"database": dbStatus},
	}

	if !status.healthy {
		render.Fail(c, http.StatusServiceUnavailable, response)
		return
	}

	render.OK(c, response)
}

func (h *HealthCheckHandler) Ping(c *gin.Context) {
	render.OK(c, gin.H{"message": "pong"})
}
