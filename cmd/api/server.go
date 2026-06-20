package main

import (
	"github.com/Milua25/go-job-application-tracker/internal/middlewares"
	"github.com/Milua25/go-job-application-tracker/internal/routers"

	"github.com/gin-gonic/gin"
)

func main() {

	// Set the router as the default one shipped with Gin
	router := gin.New()

	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	// recover from any panics and writes a 500 if there was one.
	hppOptions := middlewares.HPPOptions{
		CheckBody:        true,
		CheckQuery:       true,
		CheckContentType: "application/x-www-form-urlencoded",
		Whitelist:        []string{"tags", "categories"},
	}

	allowedMiddlewares := []gin.HandlerFunc{
		gin.Logger(),
		gin.Recovery(),
		middlewares.RequireOrigin(),
		middlewares.Cors(),
		//middlewares.ResponseTime(),
		middlewares.SecurityHeaders(),
		middlewares.SessionCSRFProtection(),
		middlewares.CSRFProtection(),
		middlewares.CSRFTokenHeader(),
		middlewares.Compression(),
		hppOptions.Hpp(),
		middlewares.NewRateLimiter(10).Limit(),
	}

	router.Use(allowedMiddlewares...)

	routers.APIv1Routes(router)

	router.Run() // listens on 0.0.0.0:8080 by default
}
