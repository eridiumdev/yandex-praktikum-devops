package main

import (
	"context"
	_http "eridiumdev/yandex-praktikum-go-devops/cmd/server/http"
	"eridiumdev/yandex-praktikum-go-devops/cmd/server/http/handlers"
	"eridiumdev/yandex-praktikum-go-devops/cmd/server/http/routers"
	"eridiumdev/yandex-praktikum-go-devops/cmd/server/rendering"
	"eridiumdev/yandex-praktikum-go-devops/cmd/server/storage"
	"eridiumdev/yandex-praktikum-go-devops/internal/logger"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	LogLevel = logger.LevelInfo

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

	// Init storage
	inMemStorage := storage.NewInMemStorage()

	// Init rendering engines
	htmlEngine := rendering.NewHTMLEngine("rendering/templates/html")

	// Init router
	router := routers.NewChiRouter()

	// Init handlers
	_ = handlers.NewMetricsHandler(router, inMemStorage, htmlEngine)

	// Init HTTP server
	server := _http.NewServer(router, HTTPHost, HTTPPort)

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
