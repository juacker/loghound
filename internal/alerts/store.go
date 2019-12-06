package alerts

import (
	"sync"
	"time"
)

// metricStore used in the alerting system
type metricStore struct {
	sync.Mutex
	points   []datapoint
	interval int64
	sum      int
}

// push adds a new item to the store
func (m *metricStore) push(datapoint datapoint) {
	m.Lock()
	defer m.Unlock()

	m.points = append(m.points, datapoint)

	m.sum += datapoint.Value
}

// stats return actual store stats, previously expiring out of interval datapoints
func (m *metricStore) stats() int {
	m.Lock()
	defer m.Unlock()

	limit := time.Now().Unix() - m.interval

	for len(m.points) > 0 {
		if m.points[0].Timestamp >= limit {
			break
		}

		m.sum -= m.points[0].Value
		m.points = m.points[1:]
	}

	return m.sum
}

// mean returns the mean value of elements in store
func (m *metricStore) mean() float64 {
	sum := m.stats()

	return float64(sum) / float64(m.interval)
}

// Datapoint represents a datapoint
type datapoint struct {
	Timestamp int64
	Value     int
}
