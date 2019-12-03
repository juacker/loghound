package message

// StatMessage struct
type StatMessage struct {
	Message
	Stats map[string]int `json:"stats"`
	Init  int64          `json:"init"`
	End   int64          `json:"end"`
}

// IsValid check if message has the right type
func (m *StatMessage) IsValid() bool {
	return m.Message.Type == TypeStat
}

// NewStatMessage returns a new StatMessage
func NewStatMessage(stats map[string]int, init, end int64) *StatMessage {
	return &StatMessage{
		Message: Message{TypeStat},
		Stats:   stats,
		Init:    init,
		End:     end,
	}
}
