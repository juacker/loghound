package alerts

import (
	"log"
	"sync"
	"time"
)

// metricStore used in the alerting system
type metricStore struct {
	sync.Mutex
	points   []Datapoint
	interval int64
	sum      int
}

// Push adds a new item to the store
func (m *metricStore) Push(datapoint Datapoint) {
	m.Lock()
	defer m.Unlock()

	m.points = append(m.points, datapoint)

	m.sum += datapoint.Value
	log.Println("POINTS ", len(m.points))
}

// Stats return actual store stats, previously expiring out of interval datapoints
func (m *metricStore) Stats() int {
	m.Lock()
	defer m.Unlock()

	limit := time.Now().Unix() - m.interval

	for len(m.points) > 0 {
		if m.points[0].Timestamp >= limit {
			break
		}
		log.Println("EXPIRING POINT", m.points[0].Timestamp, limit)

		m.sum -= m.points[0].Value
		m.points = m.points[1:]
	}

	return m.sum
}

// Mean returns the mean value of elements in store
func (m *metricStore) Mean() float64 {
	sum := m.Stats()
	log.Println("alerts: MEAN ", sum, m.interval)

	return float64(sum) / float64(m.interval)
}

// Datapoint represents a datapoint
type Datapoint struct {
	Timestamp int64
	Value     int
}
