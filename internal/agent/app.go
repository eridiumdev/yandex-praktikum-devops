package agent

import (
	"context"
	"time"

	"eridiumdev/yandex-praktikum-go-devops/config"
	"eridiumdev/yandex-praktikum-go-devops/internal/common/helpers"
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/buffering"
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/collectors"
	delivery "eridiumdev/yandex-praktikum-go-devops/internal/metrics/delivery/http"
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/exporters"
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/hash"
)

type app struct {
	*agent
}

func NewApp(ctx context.Context, cfg *config.AgentConfig) (*app, error) {
	// Init buffer for metrics
	metricsBuffer := buffering.NewInMemBuffer()

	// Init agent
	agent := newAgent(cfg, metricsBuffer)

	// Init collectors
	runtimeCollector := collectors.NewRuntimeCollector("runtime")
	pollCountCollector := collectors.NewPollCountCollector("poll-count")
	randomCollector, err := collectors.NewRandomCollector("random", cfg.RandomCollector)
	if err != nil {
		return nil, helpers.WrapErr(err, agent, "cannot start random collector")
	}

	// Provide collectors to agent
	agent.AddCollector(runtimeCollector)
	agent.AddCollector(pollCountCollector)
	agent.AddCollector(randomCollector)

	// Init auxiliary components
	hasher := hash.NewHasher(cfg.HashKey)
	requestResponseFactory := delivery.NewRequestResponseFactory(hasher)

	// Init exporters
	httpExporter := exporters.NewHTTPExporter("http", requestResponseFactory, cfg.HTTPExporter)
	agent.AddExporter(httpExporter)

	return &app{
		agent: agent,
	}, nil
}

func (a *app) Run(ctx context.Context) {
	// Wait one collectInterval before running first export
	time.AfterFunc(a.cfg.CollectInterval, func() {
		a.StartExporting(ctx)
	})
	// Start collecting (endless cycle)
	a.StartCollecting(ctx)
}

func (a *app) Stop(ctx context.Context) {
	// Wait for collectors to finish their job
	// (= try to reserve all available worker threads)
	for _, col := range a.collectors {
		for i := 0; i < col.MaxThreads(); i++ {
			col.Reserve(ctx)
		}
	}
	// Wait for exporters to finish their job
	// (= try to reserve all available worker threads)
	for _, exp := range a.exporters {
		for i := 0; i < exp.MaxThreads(); i++ {
			exp.Reserve(ctx)
		}
	}
}
