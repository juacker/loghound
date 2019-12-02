package message

// Type defines the type of a message
type Type int

// messageTypes
const (
	TypeCLF Type = iota
	TypeStat
	TypeAlert
)

// Message to map the messages sent by the modules
type Message struct {
	Type Type `json:"type"`
}
