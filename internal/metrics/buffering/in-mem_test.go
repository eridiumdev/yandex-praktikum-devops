package buffering

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/domain"
)

func TestBuffer(t *testing.T) {
	tests := []struct {
		name string
		have map[string]domain.Metric
		add  []domain.Metric
		want map[string]domain.Metric
	}{
		{
			name: "add counter to empty buffer",
			have: map[string]domain.Metric{},
			add: []domain.Metric{
				domain.NewCounter(domain.PollCount, 10),
			},
			want: map[string]domain.Metric{
				domain.PollCount: domain.NewCounter(domain.PollCount, 10),
			},
		},
		{
			name: "add gauge to empty buffer",
			have: map[string]domain.Metric{},
			add: []domain.Metric{
				domain.NewGauge(domain.Alloc, 10.333),
			},
			want: map[string]domain.Metric{
				domain.Alloc: domain.NewGauge(domain.Alloc, 10.333),
			},
		},
		{
			name: "update counter",
			have: map[string]domain.Metric{
				domain.PollCount: domain.NewCounter(domain.PollCount, 10),
			},
			add: []domain.Metric{
				domain.NewCounter(domain.PollCount, 20),
			},
			want: map[string]domain.Metric{
				domain.PollCount: domain.NewCounter(domain.PollCount, 30),
			},
		},
		{
			name: "update gauge",
			have: map[string]domain.Metric{
				domain.Alloc: domain.NewGauge(domain.Alloc, 10.333),
			},
			add: []domain.Metric{
				domain.NewGauge(domain.Alloc, 20.555),
			},
			want: map[string]domain.Metric{
				domain.Alloc: domain.NewGauge(domain.Alloc, 20.555),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := inMemBuffer{
				buffer: tt.have,
				mutex:  &sync.RWMutex{},
			}
			buf.Buffer(tt.add)
			assert.EqualValues(t, tt.want, buf.buffer)
		})
	}
}

func TestBufferWithRaceCondition(t *testing.T) {
	buffer := NewInMemBuffer()
	metric := domain.NewCounter(domain.PollCount, 1)

	done := make(chan int)
	for i := 0; i < 1000; i++ {
		go func() {
			buffer.Buffer([]domain.Metric{metric})
			done <- 1
		}()
	}
	threadsDone := 0
	for range done {
		threadsDone++
		if threadsDone == 1000 {
			break
		}
	}
	result := buffer.Retrieve()
	assert.Equal(t, domain.Counter(1000), result[0].Value())
}

func TestRetrieve(t *testing.T) {
	tests := []struct {
		name   string
		buffer map[string]domain.Metric
		want   []domain.Metric
	}{
		{
			name:   "retrieve from empty buffer",
			buffer: map[string]domain.Metric{},
			want:   []domain.Metric{},
		},
		{
			name: "retrieve from non-empty buffer",
			buffer: map[string]domain.Metric{
				domain.PollCount: domain.NewCounter(domain.PollCount, 10),
				domain.Alloc:     domain.NewGauge(domain.Alloc, 10.333),
			},
			want: []domain.Metric{
				domain.NewCounter(domain.PollCount, 10),
				domain.NewGauge(domain.Alloc, 10.333),
			},
		},
		{
			name: "retrieve from non-empty buffer, different order",
			buffer: map[string]domain.Metric{
				domain.Alloc:     domain.NewGauge(domain.Alloc, 10.333),
				domain.PollCount: domain.NewCounter(domain.PollCount, 10),
			},
			want: []domain.Metric{
				domain.NewCounter(domain.PollCount, 10),
				domain.NewGauge(domain.Alloc, 10.333),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := inMemBuffer{
				buffer: tt.buffer,
				mutex:  &sync.RWMutex{},
			}
			list := buf.Retrieve()
			assert.ElementsMatch(t, tt.want, list)
		})
	}
}

func TestFlush(t *testing.T) {
	tests := []struct {
		name   string
		buffer map[string]domain.Metric
		want   map[string]domain.Metric
	}{
		{
			name:   "flush empty buffer",
			buffer: map[string]domain.Metric{},
			want:   map[string]domain.Metric{},
		},
		{
			name: "flush non-empty buffer",
			buffer: map[string]domain.Metric{
				domain.PollCount: domain.NewCounter(domain.PollCount, 10),
				domain.Alloc:     domain.NewGauge(domain.Alloc, 10.333),
			},
			want: map[string]domain.Metric{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := inMemBuffer{
				buffer: tt.buffer,
				mutex:  &sync.RWMutex{},
			}
			buf.Flush()
			assert.EqualValues(t, tt.want, buf.buffer)
		})
	}
}
