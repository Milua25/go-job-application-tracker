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

	"github.com/Milua25/go-job-application-tracker/internal/auth"
	"github.com/Milua25/go-job-application-tracker/internal/config"
	"github.com/Milua25/go-job-application-tracker/internal/healthcheck"
	"github.com/Milua25/go-job-application-tracker/internal/middleware"
	"github.com/Milua25/go-job-application-tracker/internal/repository/sqlconnect"
	"github.com/Milua25/go-job-application-tracker/internal/routers"
	"github.com/Milua25/go-job-application-tracker/internal/token"
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
	tokenMaker, err := token.NewJWTMaker(
		a.cfg.JWT.SecretKey,
		a.cfg.JWT.Issuer,
		a.cfg.JWT.ExpiresIn,
		a.cfg.JWT.RefreshExpiresIn,
	)

	if err != nil {
		return fmt.Errorf("invalid JWT config: %w", err)
	}

	userHandler := user.NewUserHandler(a.pgStore.User)
	healthcheckHandler := healthcheck.NewHealthCheckHandler(a.sqlDB)
	authHandler := auth.NewAuthHandler(a.pgStore.User, a.pgStore.Session, tokenMaker)

	router := gin.New()
	router.Use(a.middleware...)

	routers.RegisterV1Routes(router, middleware.AuthMiddleware(tokenMaker), middleware.GetAdminMiddleware(tokenMaker), userHandler, healthcheckHandler, authHandler)

	addr := fmt.Sprintf("%s:%s", a.cfg.Server.Addr, a.cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  a.cfg.Server.ReadTimeout,
		WriteTimeout: a.cfg.Server.WriteTimeout,
		IdleTimeout:  a.cfg.Server.IdleTimeout,
	}
	serverErrors := make(chan error, 1)
	go func() {
		slog.Info("server listening", "addr", addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", "error", err)
			serverErrors <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
	case err := <-serverErrors:
		return err
	}

	slog.Info("shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), a.cfg.Server.ShutdownTimeout)
	defer cancel()

	return srv.Shutdown(ctx)
}
