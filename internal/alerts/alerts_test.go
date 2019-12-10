package alerts

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/juacker/loghound/internal/broker"
	"github.com/juacker/loghound/internal/message"
	"github.com/juacker/loghound/internal/testutils"
	tassert "github.com/stretchr/testify/assert"
)

func TestInternalAlerts(t *testing.T) {

	assert := tassert.New(t)
	t.Helper()

	link := testutils.Link{
		T: t,
	}

	monitor := &metricMonitor{
		broker: &link,
	}

	// processMessage invalid message
	t.Run("processMessage - fail unmarshaling", func(t *testing.T) {
		err := monitor.processMessage([]byte("msg"))
		assert.Equal("invalid character 'm' looking for beginning of value", err.Error())
	})

	// processMessage invalid message
	t.Run("processMessage - invalid message", func(t *testing.T) {
		msg := message.Message{
			Type: message.TypeCLF,
		}

		payload, err := json.Marshal(msg)
		assert.Nil(err, "err nil")
		assert.Equal(fmt.Errorf("invalid message"), monitor.processMessage(payload))
	})

	// processMessage success
	t.Run("processMessage - success - empty", func(t *testing.T) {
		msg := message.NewStatMessage(make(map[string]int), 0, 0)

		payload, err := json.Marshal(msg)
		assert.Nil(err, "err nil")
		assert.Equal(nil, monitor.processMessage(payload))
	})

	// processMessage success not my metric
	t.Run("processMessage - success - my metric", func(t *testing.T) {
		monitor.metric = "my.metric"
		monitor.store = &metricStore{
			points:   make([]datapoint, 0),
			interval: 1,
		}

		metrics := make(map[string]int)
		metrics["another.metric"] = 42

		msg := message.NewStatMessage(metrics, 0, time.Now().Unix())

		payload, err := json.Marshal(msg)
		assert.Nil(err, "err nil")
		assert.Equal(nil, monitor.processMessage(payload), "err nil")
		assert.Equal(0, monitor.store.sum, "expected store sum")
	})

	// processMessage success my metric
	t.Run("processMessage - success - my metric", func(t *testing.T) {
		monitor.metric = "my.metric"
		monitor.store = &metricStore{
			points:   make([]datapoint, 0),
			interval: 1,
		}

		metrics := make(map[string]int)
		metrics["my.metric"] = 42

		msg := message.NewStatMessage(metrics, 0, time.Now().Unix())

		payload, err := json.Marshal(msg)
		assert.Nil(err, "err nil")
		assert.Equal(nil, monitor.processMessage(payload), "err nil")
		assert.Equal(42, monitor.store.sum, "expected store sum")
	})

	// checkAlert success nothing to do
	t.Run("checkAlert - success - nothing to do", func(t *testing.T) {
		link.Reset()

		monitor.store = &metricStore{
			interval: 1,
		}

		assert.Nil(nil, monitor.checkAlert(), "err nil")
		assert.Equal(0, link.SendCount, "no message sent to broker")
	})

	// checkAlert success threshold raised
	t.Run("checkAlert - success - threshold raised", func(t *testing.T) {
		link.Reset()
		assert.Equal(0, link.SendCount, "initial state for link")

		monitor.threshold = 1
		monitor.metric = "my.metric"
		monitor.store = &metricStore{
			sum:      10,
			interval: 1,
		}

		topic := broker.TopicAlert
		text := fmt.Sprintf("High traffic generated an alert - hits = {%.2f}, triggered at {%v}", float64(10), time.Now().Truncate(time.Second))

		link.ExpectedSentTopic = &topic
		link.ExpectedSentMsg = message.NewAlertMessage(monitor.metric, text, message.SeverityMax)

		assert.Nil(nil, monitor.checkAlert(), "err nil")
		assert.Equal(1, link.SendCount, "message sent to broker")
	})

	// checkAlert success threshold cancelled
	t.Run("checkAlert - success - threshold cancelled", func(t *testing.T) {
		link.Reset()
		assert.Equal(0, link.SendCount, "initial state for link")

		monitor.threshold = 1
		monitor.metric = "my.metric"
		monitor.store = &metricStore{
			sum:      0,
			interval: 1,
		}

		topic := broker.TopicAlert
		text := fmt.Sprintf("High traffic alert CANCELED - hits = {%.2f}, at {%v}", 0.0, time.Now().Truncate(time.Second))

		link.ExpectedSentTopic = &topic
		link.ExpectedSentMsg = message.NewAlertMessage(monitor.metric, text, message.SeverityCanceled)

		assert.Nil(nil, monitor.checkAlert(), "err nil")
		assert.Equal(1, link.SendCount, "message sent to broker")
	})
}
