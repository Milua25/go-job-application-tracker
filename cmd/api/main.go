package main

import (
	"time"

	"github.com/Milua25/go-job-application-tracker/internal/config"
	"github.com/Milua25/go-job-application-tracker/internal/middleware"
	"github.com/Milua25/go-job-application-tracker/internal/repository/sqlconnect"

	"log/slog"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		return
	}

	hppOptions := middleware.HPPOptions{
		CheckBody:        true,
		CheckQuery:       true,
		CheckContentType: "application/x-www-form-urlencoded",
		Whitelist:        []string{"tags", "categories"},
	}

	slogJsonHandler := slog.NewJSONHandler(gin.DefaultWriter, nil)
	sloggerHandler := slog.New(slog.NewMultiHandler(slogJsonHandler))

	allowedMiddlewares := []gin.HandlerFunc{
		gin.Logger(),
		gin.Recovery(),
		middleware.RequestIDMiddleware(),
		middleware.SlogMiddleware(sloggerHandler),
		middleware.RequireOrigin(),
		middleware.Cors(),
		middleware.SecurityHeaders(),
		middleware.SessionCSRFProtection(cfg.Security.SessionSecret),
		middleware.CSRFProtection(cfg.Security.CSRFSecret),
		middleware.CSRFTokenHeader(),
		middleware.Compression(),
		hppOptions.Hpp(),
		middleware.NewRateLimiter(10).Limit(),
	}

	gormDB, err := sqlconnect.ConnectToPgDB(cfg.DB.DSN())
	if err != nil {
		slog.Error("Failed to connect to the database", "error", err)
		return
	}
	sqlDB, err := gormDB.DB()
	if err != nil {
		slog.Error("Failed to get sql.DB from gorm.DB", "error", err)
		return
	}
	sqlDB.SetMaxIdleConns(cfg.DB.MinIdleConns)
	sqlDB.SetMaxOpenConns(cfg.DB.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.DB.ConnMaxLifetime) * time.Second)
	sqlDB.SetConnMaxIdleTime(time.Duration(cfg.DB.DefaultTimeoutDuration) * time.Second)

	if err := applyDBMigrations(sqlDB); err != nil {
		slog.Error("Failed to apply database migrations", "error", err)
		return
	}

	storePostgres := sqlconnect.NewPostgresStore(gormDB)

	// Create the app instance
	app := &app{
		cfg:        cfg,
		middleware: allowedMiddlewares,
		pgStore:    storePostgres,
		sqlDB:      sqlDB,
	}

	if err := app.run(); err != nil {
		slog.Error("Failed to run the app", "error", err)
	}
}
