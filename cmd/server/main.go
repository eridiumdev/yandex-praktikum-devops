package main

import (
	"context"
	_http "eridiumdev/yandex-praktikum-go-devops/cmd/server/http"
	"eridiumdev/yandex-praktikum-go-devops/cmd/server/http/handlers"
	"eridiumdev/yandex-praktikum-go-devops/cmd/server/storage"
	"eridiumdev/yandex-praktikum-go-devops/internal/logger"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const (
	LogLevel = logger.LevelDebug

	HTTPHost = "127.0.0.1"
	HTTPPort = 8080

	ShutdownTimeout = 3 * time.Second
)

func main() {
	// Init context
	ctx := context.Background()

	// Init custom logger
	logger.Init(LogLevel)
	logger.Infof("Logger started")

	// Init HTTP server
	server := _http.NewServer(HTTPHost, HTTPPort)

	// Init storage
	inMemStorage := storage.NewInMemStorage()

	// Init handlers
	metricsHandler := handlers.NewMetricsHandler(inMemStorage)

	// Connect handlers to server
	server.AddHandler("/update/", http.MethodPost, metricsHandler.Update)

	// Start server
	logger.Infof("Starting HTTP server on %s:%d", HTTPHost, HTTPPort)
	go server.Start()

	// Handle OS signals for graceful shutdown
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Kill)
	sig := <-quit
	logger.Infof("OS signal received: %s", sig)

	// Allow some time for collectors/exporters to finish their job
	time.AfterFunc(ShutdownTimeout, func() {
		logger.Fatalf("Server force-stopped (shutdown timeout)")
	})

	// Call cancel function and stop the server
	server.Stop(ctx)
	logger.Infof("Server stopped")
}
