package healthcheck

import (
	"context"
	"database/sql"
	"log/slog"
	"time"
)

type healthCheckService struct {
	db *sql.DB
}

func newHealthCheckService(db *sql.DB) *healthCheckService {
	return &healthCheckService{db: db}
}

type healthStatus struct {
	healthy bool
	dbUp    bool
}

func (s *healthCheckService) check(ctx context.Context) healthStatus {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	dbUp := true
	if err := s.db.PingContext(ctx); err != nil {
		slog.Error("database health check failed", "error", err)
		dbUp = false
	}

	return healthStatus{healthy: dbUp, dbUp: dbUp}
}
