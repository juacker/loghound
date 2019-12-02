package stats

import (
	"sync"
)

type cache struct {
	sync.Mutex
	metrics map[string]int
}

func (c *cache) Increment(metric string, value int) {
	c.Lock()
	defer c.Unlock()

	c.metrics[metric] += value
}

func (c *cache) Stats() map[string]int {
	c.Lock()
	defer c.Unlock()

	stats := make(map[string]int)

	for k, v := range c.metrics {
		stats[k] = v
		delete(c.metrics, k)
	}

	return stats
}
