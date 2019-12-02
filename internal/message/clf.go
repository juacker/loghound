package message

import (
	"github.com/juacker/loghound/pkg/clf"
)

// CLFMessage struct
type CLFMessage struct {
	Message
	clf.Clf
}

// IsValid check if message has the right type
func (m *CLFMessage) IsValid() bool {
	return m.Message.Type == TypeCLF
}

// NewCLFMessage returns a new CLFMessage
func NewCLFMessage(m *clf.Clf) *CLFMessage {
	return &CLFMessage{
		Message{TypeCLF},
		*m,
	}
}
