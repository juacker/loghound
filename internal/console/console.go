package console

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/juacker/loghound/internal/broker"
	"github.com/juacker/loghound/internal/message"
	"github.com/juacker/loghound/pkg/clf"
)

type console struct {
	broker    *broker.Connection
	dashboard *clf.Dashboard
}

func (c *console) loop() {
	log.Println("console: initializing console")

	if err := ui.Init(); err != nil {
		log.Fatalf("console: failed to initialize termui: %v", err)
	}
	defer ui.Close()

	width, height := ui.TerminalDimensions()
	interval := int64(600)
	c.dashboard = clf.NewDashboard(width, height, interval)

	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(time.Second)
LOOP:
	for {
		select {
		case payload := <-c.broker.Receive():
			log.Println("console: new message received")
			err := c.processMessage(payload)
			if err != nil {
				log.Println("console: failed processing message: ", err)
			}
		case <-ticker.C:
			c.dashboard.Render()
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				log.Println("console: 'q' or '<C-c>' key pressed, exiting")
				break LOOP
			case "<Resize>":
				payload := e.Payload.(ui.Resize)
				c.dashboard.Resize(payload.Width, payload.Height)
			}
		}
	}
}

func (c *console) processMessage(payload []byte) error {
	var msg message.Message
	err := json.Unmarshal(payload, &msg)
	if err != nil {
		return err
	}

	switch t := msg.Type; t {
	case message.TypeAlert:
		var alertMsg message.AlertMessage
		err := json.Unmarshal(payload, &alertMsg)
		if err != nil {
			return err
		}
		return c.processAlertMessage(&alertMsg)
	case message.TypeStat:
		var statMsg message.StatMessage
		err := json.Unmarshal(payload, &statMsg)
		if err != nil {
			return err
		}
		return c.processStatMessage(&statMsg)
	default:
		log.Println("console: invalid message received")
		return nil
	}
}

func (c *console) processAlertMessage(msg *message.AlertMessage) error {
	if !msg.IsValid() {
		return fmt.Errorf("invalid message")
	}

	c.dashboard.Message(time.Now(), msg.Text)
	return nil
}

func (c *console) processStatMessage(msg *message.StatMessage) error {
	if !msg.IsValid() {
		return fmt.Errorf("invalid message")
	}

	for metric, value := range msg.Stats {
		c.dashboard.AddPoint(metric, msg.End, float64(value))
	}

	return nil
}

// Run starts console
func Run() {
	conn, err := broker.NewConnection(broker.TopicStat, broker.TopicAlert)
	if err != nil {
		log.Fatal("console: failed opening broker connection ", err)
	}

	console := &console{
		broker: conn,
	}

	console.loop()
}
