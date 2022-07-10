package buffering

import (
	"sync"

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
			b.buffer[metric.Name()].Update(metric.Value())
		} else {
			b.buffer[metric.Name()] = metric
		}
	}
}

func (b *inMemBuffer) Retrieve() []domain.Metric {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	var result []domain.Metric

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
