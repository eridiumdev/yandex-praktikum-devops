package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"eridiumdev/yandex-praktikum-go-devops/config"
	"eridiumdev/yandex-praktikum-go-devops/internal/common/logger"
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/buffering"
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/executors/collectors"
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/executors/exporters"
)

const (
	LogLevel = logger.LevelInfo
	LogMode  = logger.ModeDevelopment

	LogExporter  = 0x01
	HTTPExporter = 0x02

	// ExportersEnabled = LogExporter
	ExportersEnabled = HTTPExporter
)

func main() {
	// Init logger
	ctx := logger.Init(context.Background(), LogLevel, LogMode)
	logger.New(ctx).Infof("Logger started")

	// Get context with cancel func for graceful shutdown
	ctx, cancel := context.WithCancel(ctx)

	// Init config
	cfg, err := config.LoadAgentConfig(config.FromEnv)
	if err != nil {
		logger.New(ctx).Fatalf("Cannot load config: %s", err.Error())
	}

	// Init buffer for metrics
	metricsBuffer := buffering.NewInMemBuffer()

	// Init agent
	agent := NewAgent(cfg, metricsBuffer)

	// Init collectors
	runtimeCollector := collectors.NewRuntimeCollector("runtime")
	pollCountCollector := collectors.NewPollCountCollector("poll-count")
	randomCollector, err := collectors.NewRandomCollector("random", cfg.RandomExporter)
	if err != nil {
		logger.New(ctx).Fatalf("Cannot start random collector: %s", err.Error())
	}

	// Provide collectors to agent
	agent.AddCollector(runtimeCollector)
	agent.AddCollector(pollCountCollector)
	agent.AddCollector(randomCollector)

	// Init exporters
	if exporterEnabled(LogExporter) {
		// LogExporter is used for debug purposes
		logExporter := exporters.NewLogExporter("log")
		agent.AddExporter(logExporter)
	}
	if exporterEnabled(HTTPExporter) {
		// HTTPExporter is the main exporter for Yandex-Practicum tasks
		httpExporter := exporters.NewHTTPExporter("http", cfg.HTTPExporter)
		agent.AddExporter(httpExporter)
	}

	// Start agent
	go agent.StartCollecting(ctx)
	// Wait one collectInterval before running first export
	time.AfterFunc(cfg.CollectInterval, func() {
		agent.StartExporting(ctx)
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
	agent.Stop()
	logger.New(ctx).Infof("Agent stopped")
}

func exporterEnabled(exporter int) bool {
	return ExportersEnabled&exporter == exporter
}
