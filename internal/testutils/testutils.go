package testutils

import (
	"encoding/json"
	"testing"
	"time"

	tassert "github.com/stretchr/testify/assert"
)

func testInternalAlerts(t *testing.T) {

	t.Run("prefix", func(t *testing.T) {

	})
}

// Link satisfies broker.Link interface
type Link struct {
	T                   *testing.T
	c                   chan []byte
	SendCount           int
	ReceiveCount        int
	ExpectedSentTopic   *int
	ExpectedSentMsg     interface{}
	ExpectedReceivedMsg interface{}
}

// Reset reset struct fields
func (l *Link) Reset() {
	l.c = make(chan []byte)
	l.SendCount = 0
	l.ReceiveCount = 0
	l.ExpectedSentTopic = nil
	l.ExpectedSentMsg = nil
	l.ExpectedReceivedMsg = nil
}

// Send validates sent messages to the link and increases call counter
func (l *Link) Send(topic int, msg interface{}) error {
	assert := tassert.New(l.T)

	l.SendCount++
	assert.NotNil(l.ExpectedSentTopic, "we expect a message to be sent")
	assert.Equal(*l.ExpectedSentTopic, topic, "topic is the expected")
	assert.Equal(l.ExpectedSentMsg, msg, "message is the expected")
	return nil
}

// Receive increases call counter and returns the expected message
func (l *Link) Receive() <-chan []byte {
	assert := tassert.New(l.T)

	l.ReceiveCount++

	go func() {
		time.Sleep(100 * time.Millisecond)
		payload, err := json.Marshal(l.ExpectedReceivedMsg)
		assert.NotNil(err, "valid message expected")
		l.c <- payload
	}()

	return l.c
}
