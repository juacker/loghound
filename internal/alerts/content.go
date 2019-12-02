package alerts

// SortableContent describe a content that can be added to a sorted cache
type SortableContent interface {
	SortedValue() int64
}

// Datapoint represents a datapoint
type Datapoint struct {
	Timestamp int64
	Metric    string
	Value     int
}

// SortedValue returns the metric timestamp
func (d Datapoint) SortedValue() int64 {
	return d.Timestamp
}
