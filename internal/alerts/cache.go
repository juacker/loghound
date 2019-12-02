package alerts

import (
	"container/heap"
	"sync"
)

// cache used in the alerting system
type cache struct {
	sync.Mutex
	elements   Elements
	totalValue int
}

// Add adds a new item to the cache
func (c *cache) Add(content SortableContent) {
	c.Lock()
	defer c.Unlock()

	element := &Element{
		value:   content.SortedValue(),
		Content: content,
	}

	heap.Push(&c.elements, element)
	c.totalValue += content.(Datapoint).Value
}

// Purge removes elements less than value
func (c *cache) Purge(value int64) {
	c.Lock()
	defer c.Unlock()

	for c.elements.Len() > 0 {
		if c.elements[0].value > value {
			break
		}

		content := heap.Pop(&c.elements)
		c.totalValue -= content.(Datapoint).Value
	}
}

// Mean returns the mean value of elements in cache
func (c *cache) Mean() float64 {
	c.Lock()
	defer c.Unlock()

	return float64(c.totalValue) / float64(c.elements.Len())
}

// Elements is a min-Heap of Element
type Elements []*Element

func (e Elements) Len() int           { return len(e) }
func (e Elements) Less(i, j int) bool { return e[i].value < e[j].value }
func (e Elements) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }

// Push adds an element to the heap
func (e *Elements) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	n := len(*e)
	item := x.(*Element)
	item.index = n
	*e = append(*e, item)
}

// Pop removes the min element from the heap
func (e *Elements) Pop() interface{} {
	old := *e
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*e = old[0 : n-1]
	return item
}

// Element is an element in the heap
type Element struct {
	index   int
	value   int64
	Content interface{}
}
