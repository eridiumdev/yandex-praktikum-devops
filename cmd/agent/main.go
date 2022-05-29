package main

import (
	"eridiumdev/yandex-praktikum-go-devops/cmd/agent/collectors"
	"eridiumdev/yandex-praktikum-go-devops/cmd/agent/exporters"
	"eridiumdev/yandex-praktikum-go-devops/internal/logger"
	"time"
)

const (
	LogExporter  = 0x01
	HttpExporter = 0x02
)

const (
	LogLevel = logger.LevelInfo

	//ExportersEnabled = LogExporter
	ExportersEnabled = HttpExporter

	CollectInterval = 2 * time.Second
	ExportInterval  = 10 * time.Second

	RandomValueMin = 0
	RandomValueMax = 9999

	HttpHost    = "127.0.0.1"
	HttpPort    = 8080
	HttpTimeout = 3 * time.Second
)

func main() {
	// Init custom logger
	logger.Init(LogLevel)
	logger.Infof("Logger started")

	// Init agent
	agent := NewAgent(CollectInterval, ExportInterval)

	// Init collectors
	runtimeCollector := collectors.NewRuntimeCollector("runtime")
	randomCollector, err := collectors.NewRandomCollector("random", RandomValueMin, RandomValueMax)
	if err != nil {
		logger.Fatalf("cannot start random collector: %s", err.Error())
	}

	// Provide collectors to agent
	agent.AddCollector(runtimeCollector)
	agent.AddCollector(randomCollector)

	// Init exporters
	if exporterEnabled(LogExporter) {
		// LogExporter is used for debug purposes
		logExporter := exporters.NewLogExporter("log")
		agent.AddExporter(logExporter)
	}
	if exporterEnabled(HttpExporter) {
		// HttpExporter is the main exporter for Yandex-Practicum tasks
		httpExporter := exporters.NewHttpExporter("http", HttpHost, HttpPort, HttpTimeout)
		agent.AddExporter(httpExporter)
	}

	// Start agent
	go agent.StartCollecting()
	// Wait one CollectInterval before running first export
	time.AfterFunc(CollectInterval, func() {
		agent.StartExporting()
	})

	logger.Infof("Agent started")
	agent.StartBuffering()
}

func exporterEnabled(exporter int) bool {
	return ExportersEnabled & exporter == exporter
}
