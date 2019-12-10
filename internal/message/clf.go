package message

import (
	"github.com/juacker/loghound/pkg/clf"
)

// CLFMessage struct
type CLFMessage struct {
	Message
	clf.Entry
}

// IsValid check if message has the right type
func (m *CLFMessage) IsValid() bool {
	return m.Message.Type == TypeCLF
}

// NewCLFMessage returns a new CLFMessage
func NewCLFMessage(m *clf.Entry) *CLFMessage {
	return &CLFMessage{
		Message{TypeCLF},
		*m,
	}
}
