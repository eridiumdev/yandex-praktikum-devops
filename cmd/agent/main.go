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

	// Init agent app
	agentCtx := logger.Enrich(ctx, logger.FieldComponent, "agent")
	app, err := agent.NewApp(agentCtx, cfg)
	if err != nil {
		logger.New(ctx).Err(err).Fatalf("Cannot init agent")
	}
	logger.New(ctx).Infof("Agent init successful")

	// Run the app
	go app.Run(agentCtx)

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
	app.Stop(agentCtx)
	logger.New(ctx).Infof("Agent stopped")
}
