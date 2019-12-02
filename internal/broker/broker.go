package broker

import (
	"encoding/json"
	"log"
	"sync"
)

type messageBroker struct {
	sync.Mutex
	ctl         chan bool
	wg          *sync.WaitGroup
	subscribers map[int][]chan []byte
	listener    chan []byte
}

type message struct {
	Topic   int    `json:"topic"`
	Payload []byte `json:"payload"`
}

var broker *messageBroker

func (b *messageBroker) listen() (err error) {
	log.Println("broker: listening messages")

LOOP:
	for {
		select {
		case payload := <-b.listener:
			log.Println("broker: new message received ", payload)
			err = b.broadcast(payload)
			if err != nil {
				log.Println("broker: error sending broadcast message", err)
			}
		case <-b.ctl:
			log.Println("broker: ctl signal received, exiting")
			break LOOP
		}
	}

	// close subscribers channels
	for _, v := range b.subscribers {
		for _, c := range v {
			if c != nil {
				close(c)
			}
		}
	}

	b.wg.Done()
	return err
}

func (b *messageBroker) broadcast(payload []byte) error {
	var msg message
	err := json.Unmarshal(payload, &msg)

	if err != nil {
		return err
	}

	b.Lock()
	defer b.Unlock()

	log.Println("Checking subscribers for topic ", msg.Topic)
	for _, subscriber := range b.subscribers[msg.Topic] {
		log.Println("sending message to subscriber ", subscriber, msg.Topic)
		subscriber <- msg.Payload
	}

	return nil
}

func (b *messageBroker) subscribe(topic ...int) (chan []byte, error) {
	b.Lock()
	defer b.Unlock()

	ch := make(chan []byte, 100)
	for _, t := range topic {
		if _, ok := b.subscribers[t]; !ok {
			b.subscribers[t] = make([]chan []byte, 0)
		}
		b.subscribers[t] = append(b.subscribers[t], ch)
	}

	return ch, nil
}

// Run starts message broker
func Run(wg *sync.WaitGroup, ctl chan bool) {
	broker.ctl = ctl
	broker.wg = wg
	broker.listen()
}

// NewConnection returns a new broker connection, the connection
// will receive messages from the subscribed topics only
func NewConnection(topic ...int) (*Connection, error) {
	readChannel, err := broker.subscribe(topic...)
	if err != nil {
		return nil, err
	}

	return &Connection{
		read:  readChannel,
		write: broker.listener,
	}, nil
}

func init() {
	broker = &messageBroker{
		subscribers: make(map[int][]chan []byte),
		listener:    make(chan []byte, 100),
	}
}
