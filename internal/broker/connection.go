package broker

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Link interface defines a link to a broker
type Link interface {
	Send(int, interface{}) error
	Receive() <-chan []byte
}

// Connection represents a broker connection. it satisfies Link interface
type Connection struct {
	read  <-chan []byte
	write chan<- []byte
}

// Send is used to send messages to the broker
func (c *Connection) Send(topic int, msg interface{}) error {
	payload, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed creating message payload")
	}

	brokerMessage := message{
		Topic:   topic,
		Payload: payload,
	}

	data, err := json.Marshal(brokerMessage)
	if err != nil {
		return errors.New("invalid message")
	}

	c.write <- data
	return nil
}

// Receive will receive messages from the broker on the subscribed topics
func (c *Connection) Receive() <-chan []byte {
	return c.read
}
