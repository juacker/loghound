package stats

import (
	"sync"
	"time"
)

type cache struct {
	sync.Mutex
	metrics map[string]int
	reset   int64
}

func (c *cache) Increment(metric string, value int) {
	c.Lock()
	defer c.Unlock()

	c.metrics[metric] += value
}

func (c *cache) Stats() (map[string]int, int64, int64) {
	c.Lock()
	defer c.Unlock()

	now := time.Now().Unix()
	begin := c.reset
	c.reset = now

	stats := make(map[string]int)

	for k, v := range c.metrics {
		stats[k] = v
		c.metrics[k] = 0
	}

	return stats, begin, now
}
