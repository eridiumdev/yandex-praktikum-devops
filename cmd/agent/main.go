package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"eridiumdev/yandex-praktikum-go-devops/config"
	"eridiumdev/yandex-praktikum-go-devops/internal/agent"
	"eridiumdev/yandex-praktikum-go-devops/internal/common/logger"
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/buffering"
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/collectors"
	delivery "eridiumdev/yandex-praktikum-go-devops/internal/metrics/delivery/http"
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/exporters"
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/hash"
)

func main() {
	// Init context
	ctx := context.Background()

	// Init config
	cfg, err := config.LoadAgentConfig()
	if err != nil {
		log.Fatalf("Cannot load config: %s", err.Error())
	}

	// Init logger and update context
	ctx = logger.InitZerolog(context.Background(), cfg.Logger)
	logger.New(ctx).Infof("Logger started")

	// Modify context with cancel func for graceful shutdown
	ctx, cancel := context.WithCancel(ctx)

	// Init buffer for metrics
	metricsBuffer := buffering.NewInMemBuffer()

	// Init agent app
	app := agent.NewAgent(cfg, metricsBuffer)

	// Init collectors
	runtimeCollector := collectors.NewRuntimeCollector("runtime")
	pollCountCollector := collectors.NewPollCountCollector("poll-count")
	randomCollector, err := collectors.NewRandomCollector("random", cfg.RandomExporter)
	if err != nil {
		logger.New(ctx).Fatalf("Cannot start random collector: %s", err.Error())
	}

	// Provide collectors to agent
	app.AddCollector(runtimeCollector)
	app.AddCollector(pollCountCollector)
	app.AddCollector(randomCollector)

	// Init auxiliary components
	hasher := hash.NewHasher(cfg.HashKey)
	requestResponseFactory := delivery.NewRequestResponseFactory(hasher)

	// Init exporters
	httpExporter := exporters.NewHTTPExporter("http", requestResponseFactory, cfg.HTTPExporter)
	app.AddExporter(httpExporter)

	// Start agent
	go app.StartCollecting(ctx)
	// Wait one collectInterval before running first export
	time.AfterFunc(cfg.CollectInterval, func() {
		app.StartExporting(ctx)
	})
	logger.New(ctx).Infof("Agent started")

	// Handle OS signals for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	logger.New(ctx).Infof("OS signal received: %s", sig)

	// Allow some time for collectors/exporters to finish their job
	time.AfterFunc(cfg.ShutdownTimeout, func() {
		logger.New(ctx).Fatalf("Agent force-stopped (shutdown timeout)")
	})

	// Call cancel function and stop the agent
	cancel()
	app.Stop(ctx)
	logger.New(ctx).Infof("Agent stopped")
}
