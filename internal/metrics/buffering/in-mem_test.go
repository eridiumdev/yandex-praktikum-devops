package buffering

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/domain"
)

func TestBuffer(t *testing.T) {
	var (
		basicCounter    = domain.NewCounter(domain.PollCount, 10)
		basicCounterUpd = domain.NewCounter(domain.PollCount, 30)
		basicGauge      = domain.NewGauge(domain.Alloc, 10.333)
		basicGaugeUpd   = domain.NewGauge(domain.Alloc, 20.555)
	)
	tests := []struct {
		name string
		have map[string]*domain.Metric
		add  []domain.Metric
		want map[string]*domain.Metric
	}{
		{
			name: "add counter to empty buffer",
			have: map[string]*domain.Metric{},
			add: []domain.Metric{
				domain.NewCounter(domain.PollCount, 10),
			},
			want: map[string]*domain.Metric{
				domain.PollCount: &basicCounter,
			},
		},
		{
			name: "add gauge to empty buffer",
			have: map[string]*domain.Metric{},
			add: []domain.Metric{
				domain.NewGauge(domain.Alloc, 10.333),
			},
			want: map[string]*domain.Metric{
				domain.Alloc: &basicGauge,
			},
		},
		{
			name: "update counter",
			have: map[string]*domain.Metric{
				domain.PollCount: &basicCounter,
			},
			add: []domain.Metric{
				domain.NewCounter(domain.PollCount, 20),
			},
			want: map[string]*domain.Metric{
				domain.PollCount: &basicCounterUpd,
			},
		},
		{
			name: "update gauge",
			have: map[string]*domain.Metric{
				domain.Alloc: &basicGauge,
			},
			add: []domain.Metric{
				domain.NewGauge(domain.Alloc, 20.555),
			},
			want: map[string]*domain.Metric{
				domain.Alloc: &basicGaugeUpd,
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

	count := 1000
	wg := sync.WaitGroup{}
	wg.Add(count)

	for i := 0; i < count; i++ {
		go func() {
			buffer.Buffer([]domain.Metric{domain.NewCounter(domain.PollCount, 1)})
			wg.Done()
		}()
	}
	wg.Wait()

	result := buffer.Retrieve()
	assert.Equal(t, domain.Counter(1000), result[0].Counter)
}

func TestRetrieve(t *testing.T) {
	var (
		basicCounter = domain.NewCounter(domain.PollCount, 10)
		basicGauge   = domain.NewGauge(domain.Alloc, 10.333)
	)
	tests := []struct {
		name   string
		buffer map[string]*domain.Metric
		want   []domain.Metric
	}{
		{
			name:   "retrieve from empty buffer",
			buffer: map[string]*domain.Metric{},
			want:   []domain.Metric{},
		},
		{
			name: "retrieve from non-empty buffer",
			buffer: map[string]*domain.Metric{
				domain.PollCount: &basicCounter,
				domain.Alloc:     &basicGauge,
			},
			want: []domain.Metric{
				domain.NewCounter(domain.PollCount, 10),
				domain.NewGauge(domain.Alloc, 10.333),
			},
		},
		{
			name: "retrieve from non-empty buffer, different order",
			buffer: map[string]*domain.Metric{
				domain.Alloc:     &basicGauge,
				domain.PollCount: &basicCounter,
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
	var (
		basicCounter = domain.NewCounter(domain.PollCount, 10)
		basicGauge   = domain.NewGauge(domain.Alloc, 10.333)
	)
	tests := []struct {
		name   string
		buffer map[string]*domain.Metric
		want   map[string]*domain.Metric
	}{
		{
			name:   "flush empty buffer",
			buffer: map[string]*domain.Metric{},
			want:   map[string]*domain.Metric{},
		},
		{
			name: "flush non-empty buffer",
			buffer: map[string]*domain.Metric{
				domain.PollCount: &basicCounter,
				domain.Alloc:     &basicGauge,
			},
			want: map[string]*domain.Metric{},
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
