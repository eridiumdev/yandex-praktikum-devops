package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/middleware"

	_http "eridiumdev/yandex-praktikum-go-devops/cmd/server/delivery/http"
	"eridiumdev/yandex-praktikum-go-devops/cmd/server/delivery/http/handlers"
	"eridiumdev/yandex-praktikum-go-devops/cmd/server/delivery/http/routers"
	"eridiumdev/yandex-praktikum-go-devops/internal/commons/logger"
	"eridiumdev/yandex-praktikum-go-devops/internal/commons/rendering"
	metricsRendering "eridiumdev/yandex-praktikum-go-devops/internal/metrics/server/rendering"
	metricsRepository "eridiumdev/yandex-praktikum-go-devops/internal/metrics/server/repository"
)

const (
	LogLevel = logger.LevelInfo
	LogMode  = logger.ModeDevelopment

	HTTPHost = "0.0.0.0"
	HTTPPort = 8080

	ShutdownTimeout = 3 * time.Second
)

func main() {
	// Init context
	ctx := context.Background()

	// Init logger
	logger.Init(LogLevel, LogMode)
	logger.Infof("Logger started")

	// Init repos
	metricsRepo := metricsRepository.NewInMemRepo()

	// Init rendering engines
	templateParser := rendering.NewHTMLTemplateParser("rendering/templates")
	metricsRenderer := metricsRendering.NewHTMLEngine(templateParser)

	// Init router
	router := routers.NewChiRouter(logger.Middleware, middleware.Recoverer)

	// Init handlers
	_ = handlers.NewMetricsHandler(router, metricsRepo, metricsRenderer)

	// Init HTTP server
	server := _http.NewServer(router, _http.ServerSettings{
		Host: HTTPHost,
		Port: HTTPPort,
	})

	// Start server
	logger.Infof("Starting HTTP server on %s:%d", HTTPHost, HTTPPort)
	go server.Start()

	// Handle OS signals for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	logger.Infof("OS signal received: %s", sig)

	// Allow some time for collectors/exporters to finish their job
	time.AfterFunc(ShutdownTimeout, func() {
		logger.Fatalf("Server force-stopped (shutdown timeout)")
	})

	// Stop the server
	server.Stop(ctx)
	logger.Infof("Server stopped")
}
