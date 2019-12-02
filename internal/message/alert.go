package message

// AlertMessage struct
type AlertMessage struct {
	Message
	Metric   string   `json:"metric"`
	Severity Severity `json:"severity"`
	Text     string   `json:"text"`
}

// IsValid check if message has the right type
func (m *AlertMessage) IsValid() bool {
	return m.Message.Type == TypeAlert
}

// Severity sets the impact of the alert
type Severity int

// Severity levels
const (
	SeverityCanceled Severity = iota
	SeverityMax
)

// NewAlertMessage returns a new AlertMessage
func NewAlertMessage(metric, text string, severity Severity) *AlertMessage {
	return &AlertMessage{
		Message:  Message{TypeAlert},
		Metric:   metric,
		Severity: severity,
		Text:     text,
	}
}
