package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Milua25/go-job-application-tracker/internal/auth"
	"github.com/Milua25/go-job-application-tracker/internal/config"
	"github.com/Milua25/go-job-application-tracker/internal/healthcheck"
	"github.com/Milua25/go-job-application-tracker/internal/middleware"
	"github.com/Milua25/go-job-application-tracker/internal/repository/sqlconnect"
	"github.com/Milua25/go-job-application-tracker/internal/routers"
	"github.com/Milua25/go-job-application-tracker/internal/tokens"
	"github.com/Milua25/go-job-application-tracker/internal/user"
	"github.com/gin-gonic/gin"
)

type app struct {
	cfg        *config.Config
	middleware []gin.HandlerFunc
	pgStore    *sqlconnect.PostgresStore
	sqlDB      *sql.DB
}

func (a *app) run() error {
	authService := tokens.NewAuthService(a.cfg.JWT.SecretKey, a.cfg.JWT.Issuer, a.cfg.JWT.ExpiresIn, a.cfg.JWT.RefreshExpiresIn)

	userHandler := user.NewUserHandler(a.pgStore.User)
	healthcheckHandler := healthcheck.NewHealthCheckHandler(a.sqlDB)
	authHandler := auth.NewAuthHandler(a.pgStore.User, authService)

	router := gin.New()
	router.Use(a.middleware...)

	routers.APIv1Routes(router, routers.Handlers{
		User:           userHandler,
		HealthCheck:    healthcheckHandler,
		Auth:           authHandler,
		AuthMiddleware: middleware.AuthMiddleware(newJWTAuthAdapter(authService)),
	})

	addr := fmt.Sprintf("%s:%s", a.cfg.Server.Addr, a.cfg.Server.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	go func() {
		slog.Info("server listening", "addr", addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return srv.Shutdown(ctx)
}
