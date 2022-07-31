package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/middleware"

	"eridiumdev/yandex-praktikum-go-devops/config"
	"eridiumdev/yandex-praktikum-go-devops/internal/common/logger"
	"eridiumdev/yandex-praktikum-go-devops/internal/common/routing"
	"eridiumdev/yandex-praktikum-go-devops/internal/common/templating"
	metricsHttpDelivery "eridiumdev/yandex-praktikum-go-devops/internal/metrics/delivery/http"
	metricsRendering "eridiumdev/yandex-praktikum-go-devops/internal/metrics/rendering"
	metricsRepository "eridiumdev/yandex-praktikum-go-devops/internal/metrics/repository"
	_metricsService "eridiumdev/yandex-praktikum-go-devops/internal/metrics/service"
)

const (
	LogLevel = logger.LevelInfo
	LogMode  = logger.ModeDevelopment
)

func main() {
	// Init logger and context
	ctx := logger.Init(context.Background(), LogLevel, LogMode)
	logger.New(ctx).Infof("Logger started")

	// Init config
	cfg, err := config.LoadServerConfig(config.FromEnv)
	if err != nil {
		logger.New(ctx).Fatalf("Cannot load config: %s", err.Error())
	}

	// Init repos
	metricsRepo := metricsRepository.NewInMemRepo()

	// Init services
	metricsService := _metricsService.NewMetricsService(metricsRepo)

	// Init rendering engines
	templateParser := templating.NewHTMLTemplateParser("web/templates")
	metricsRenderer := metricsRendering.NewHTMLEngine(templateParser)

	// Init router
	router := routing.NewChiRouter(logger.Middleware, routing.URLMiddleware, middleware.Recoverer)

	// Init handlers
	_ = metricsHttpDelivery.NewMetricsHandler(router, metricsService, metricsRenderer)

	// Init HTTP server
	server := NewServer(router.GetHandler(), cfg)

	// Start server
	logger.New(ctx).Infof("Starting HTTP server on %s", cfg.Address)
	go server.Start(ctx)

	// Handle OS signals for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	logger.New(ctx).Infof("OS signal received: %s", sig)

	// Allow some time for collectors/exporters to finish their job
	time.AfterFunc(cfg.ShutdownTimeout, func() {
		logger.New(ctx).Fatalf("Server force-stopped (shutdown timeout)")
	})

	// Stop the server
	server.Stop(ctx)
	logger.New(ctx).Infof("Server stopped")
}
