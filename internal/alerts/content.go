package alerts

// SortableContent describe a content that can be added to a sorted cache
type SortableContent interface {
	SortedValue() int64
}

// OldDatapoint represents a datapoint
type OldDatapoint struct {
	Timestamp int64
	Metric    string
	Value     int
}

// SortedValue returns the metric timestamp
func (d OldDatapoint) SortedValue() int64 {
	return d.Timestamp
}
