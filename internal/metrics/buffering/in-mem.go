package buffering

import (
	"sync"

	"eridiumdev/yandex-praktikum-go-devops/internal/commons/logger"
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/domain"
)

type inMemBuffer struct {
	buffer map[string]domain.Metric
	mutex  *sync.RWMutex
}

func NewInMemBuffer() *inMemBuffer {
	return &inMemBuffer{
		buffer: make(map[string]domain.Metric),
		mutex:  &sync.RWMutex{},
	}
}

func (b *inMemBuffer) Buffer(mtx []domain.Metric) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	for _, metric := range mtx {
		if _, ok := b.buffer[metric.Name()]; ok {
			switch metric.Type() {
			case domain.TypeCounter:
				// For counters, new value is added on top of previous value
				err := b.buffer[metric.Name()].Add(metric.Value())
				if err != nil {
					logger.Errorf("[bufferer] error when calling Add(%v) on metric %s: %s",
						metric.Value(), metric.Name(), err.Error())
				}
			case domain.TypeGauge:
				fallthrough
			default:
				// For gauges, previous value is overwritten
				err := b.buffer[metric.Name()].Set(metric.Value())
				if err != nil {
					logger.Errorf("[bufferer] error when calling Set(%v) on metric %s: %s",
						metric.Value(), metric.Name(), err.Error())
				}
			}
		} else {
			// Add copy of metric to the buffer
			b.buffer[metric.Name()] = metric.Copy()
		}
	}
}

func (b *inMemBuffer) Retrieve() []domain.Metric {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	result := make([]domain.Metric, 0)

	for _, metric := range b.buffer {
		result = append(result, metric.Copy())
	}
	return result
}

func (b *inMemBuffer) Flush() {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.buffer = make(map[string]domain.Metric)
}
