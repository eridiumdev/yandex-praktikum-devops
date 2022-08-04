package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"eridiumdev/yandex-praktikum-go-devops/config"
	"eridiumdev/yandex-praktikum-go-devops/internal/common/logger"
	"eridiumdev/yandex-praktikum-go-devops/internal/common/middleware"
	"eridiumdev/yandex-praktikum-go-devops/internal/common/routing"
	"eridiumdev/yandex-praktikum-go-devops/internal/common/templating"
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/backup"
	metricsHttpDelivery "eridiumdev/yandex-praktikum-go-devops/internal/metrics/delivery/http"
	metricsRendering "eridiumdev/yandex-praktikum-go-devops/internal/metrics/rendering"
	metricsRepository "eridiumdev/yandex-praktikum-go-devops/internal/metrics/repository"
	_metricsService "eridiumdev/yandex-praktikum-go-devops/internal/metrics/service"
	"eridiumdev/yandex-praktikum-go-devops/internal/server"
)

func main() {
	// Init context
	ctx := context.Background()

	// Init config
	cfg, err := config.LoadServerConfig()
	if err != nil {
		log.Fatalf("Cannot load config: %s", err.Error())
	}

	// Init logger and update context
	ctx = logger.Init(context.Background(), cfg.Logger)
	logger.New(ctx).Infof("Logger started")

	// Modify context with cancel func for graceful shutdown
	ctx, cancel := context.WithCancel(ctx)

	// Init repos
	inMemRepo := metricsRepository.NewInMemRepo()

	// Init backupers
	fileBackuper, err := backup.NewFileBackuper(ctx, cfg.FileBackuperPath)
	if err != nil {
		logger.New(ctx).Fatalf("Cannot init file backuper: %s", err.Error())
	}

	// Init services
	metricsService, err := _metricsService.NewMetricsService(ctx, inMemRepo, fileBackuper, cfg.Backup)
	if err != nil {
		logger.New(ctx).Fatalf("Cannot init metrics service: %s", err.Error())
	}

	// Init rendering engines
	templateParser := templating.NewHTMLTemplateParser("web/templates")
	metricsRenderer := metricsRendering.NewHTMLEngine(templateParser)

	// Init router
	router := routing.NewChiRouter(middleware.URLTrimmer)

	// Init handlers
	metricsHandler := metricsHttpDelivery.NewMetricsHandler(metricsService, metricsRenderer)
	router.AddRoute(http.MethodGet, "/", metricsHandler.List, middleware.BasicSet...)
	router.AddRoute(http.MethodPost, "/value", metricsHandler.Get, middleware.ExtendedSet...)
	router.AddRoute(http.MethodPost, "/update", metricsHandler.Update, middleware.ExtendedSet...)

	// Init HTTP server app
	app := server.NewServer(router.GetHandler(), cfg)

	// Start server
	logger.New(ctx).Infof("Starting HTTP app on %s", cfg.Address)
	go app.Start(ctx)

	// Handle OS signals for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	logger.New(ctx).Infof("OS signal received: %s", sig)

	// Allow some time for server and components to clean up
	time.AfterFunc(cfg.ShutdownTimeout, func() {
		cleanup(cancel, time.Second)
		logger.New(ctx).Fatalf("Server force-stopped (shutdown timeout)")
	})

	// Stop the server
	app.Stop(ctx)
	logger.New(ctx).Infof("Server stopped")

	// Clean-up other components, e.g. backuper
	cleanup(cancel, time.Second)
}

func cleanup(cancel context.CancelFunc, wait time.Duration) {
	cancel()
	time.Sleep(wait)
	logger.New(context.Background()).Infof("Clean-up finished")
}
