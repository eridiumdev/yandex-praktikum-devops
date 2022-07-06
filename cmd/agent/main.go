package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"eridiumdev/yandex-praktikum-go-devops/cmd/agent/collectors"
	"eridiumdev/yandex-praktikum-go-devops/cmd/agent/exporters"
	"eridiumdev/yandex-praktikum-go-devops/internal/logger"
)

const (
	LogExporter  = 0x01
	HTTPExporter = 0x02
)

const (
	LogLevel = logger.LevelInfo

	//ExportersEnabled = LogExporter
	ExportersEnabled = HTTPExporter

	CollectInterval = 2 * time.Second
	ExportInterval  = 10 * time.Second
	ShutdownTimeout = 3 * time.Second

	RandomValueMin = 0
	RandomValueMax = 9999

	HTTPHost    = "127.0.0.1"
	HTTPPort    = 8080
	HTTPTimeout = 3 * time.Second
)

func main() {
	// Init context
	ctx, cancel := context.WithCancel(context.Background())

	// Init custom logger
	logger.Init(LogLevel)
	logger.Infof("Logger started")

	// Init agent
	agent := NewAgent(CollectInterval, ExportInterval)

	// Init collectors
	runtimeCollector := collectors.NewRuntimeCollector("runtime")
	pollCountCollector := collectors.NewPollCountCollector("poll-count")
	randomCollector, err := collectors.NewRandomCollector("random", RandomValueMin, RandomValueMax)
	if err != nil {
		logger.Fatalf("Cannot start random collector: %s", err.Error())
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
		httpExporter := exporters.NewHTTPExporter("http", HTTPHost, HTTPPort, HTTPTimeout)
		agent.AddExporter(httpExporter)
	}

	// Start agent
	go agent.StartCollecting(ctx)
	// Buffering will buffer metrics before exporting
	go agent.StartBuffering(ctx)
	// Wait one CollectInterval before running first export
	time.AfterFunc(CollectInterval, func() {
		agent.StartExporting(ctx)
	})
	logger.Infof("Agent started")

	// Handle OS signals for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	logger.Infof("OS signal received: %s", sig)

	// Allow some time for collectors/exporters to finish their job
	time.AfterFunc(ShutdownTimeout, func() {
		logger.Fatalf("Agent force-stopped (shutdown timeout)")
	})

	// Call cancel function and stop the agent
	cancel()
	agent.Stop()
	logger.Infof("Agent stopped")
}

func exporterEnabled(exporter int) bool {
	return ExportersEnabled&exporter == exporter
}
