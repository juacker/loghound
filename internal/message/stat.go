package message

// StatMessage struct
type StatMessage struct {
	Message
	Stats map[string]int `json:"stats"`
}

// IsValid check if message has the right type
func (m *StatMessage) IsValid() bool {
	return m.Message.Type == TypeStat
}

// NewStatMessage returns a new StatMessage
func NewStatMessage(stats map[string]int) *StatMessage {
	return &StatMessage{
		Message: Message{TypeStat},
		Stats:   stats,
	}
}
