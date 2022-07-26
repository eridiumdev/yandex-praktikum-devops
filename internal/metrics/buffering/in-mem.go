package buffering

import (
	"sync"

	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/domain"
)

type inMemBuffer struct {
	buffer map[string]*domain.Metric
	mutex  *sync.RWMutex
}

func NewInMemBuffer() *inMemBuffer {
	return &inMemBuffer{
		buffer: make(map[string]*domain.Metric),
		mutex:  &sync.RWMutex{},
	}
}

func (b *inMemBuffer) Buffer(mtx []domain.Metric) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	for _, metric := range mtx {
		if _, ok := b.buffer[metric.Name]; ok {
			switch metric.Type {
			case domain.TypeCounter:
				// For counters, new value is added on top of previous value with AddCounter()
				b.buffer[metric.Name].Counter += metric.Counter
			case domain.TypeGauge:
				// For gauges, previous value is overwritten with SetGauge()
				b.buffer[metric.Name].Gauge = metric.Gauge
			}
		} else {
			// Add metric to the buffer
			b.buffer[metric.Name] = &metric
		}
	}
}

func (b *inMemBuffer) Retrieve() []domain.Metric {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	result := make([]domain.Metric, 0)

	for _, metric := range b.buffer {
		result = append(result, *metric)
	}
	return result
}

func (b *inMemBuffer) Flush() {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.buffer = make(map[string]*domain.Metric)
}
