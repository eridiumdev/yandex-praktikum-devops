package main

import (
	_http "eridiumdev/yandex-praktikum-go-devops/cmd/server/http"
	"eridiumdev/yandex-praktikum-go-devops/cmd/server/http/handlers"
	"eridiumdev/yandex-praktikum-go-devops/internal/logger"
	"net/http"
)

const (
	LogLevel = logger.LevelInfo

	HTTPHost = "127.0.0.1"
	HTTPPort = 8080
)

func main() {
	// Init custom logger
	logger.Init(LogLevel)
	logger.Infof("Logger started")

	// Init HTTP server
	server := _http.NewServer(HTTPHost, HTTPPort)

	// Init handlers
	metricsHandler := handlers.NewMetricsHandler()

	// Connect handlers to server
	server.AddHandler("/update/", http.MethodPost, metricsHandler.Update)

	// Start server
	logger.Infof("Starting HTTP server on %s:%d", HTTPHost, HTTPPort)
	server.Start()
}
