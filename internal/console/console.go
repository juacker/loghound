package console

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
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
			log.Println("console: new message received", payload)
			err := c.processMessage(payload)
			if err != nil {
				log.Println("console: failed processing message: ", err)
			}
		case <-ticker.C:
			c.dashboard.AddPoint("requests.total", time.Now().Unix(), rand.Float64())
			c.dashboard.AddPoint("bytes.total", time.Now().Unix(), rand.Float64())
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
	var msg message.CLFMessage
	err := json.Unmarshal(payload, &msg)
	if err != nil {
		return err
	}

	// check message is the expected
	if !msg.IsValid() {
		return fmt.Errorf("invalid message")
	}

	return c.processCLFMessage(&msg)
}

func (c *console) processCLFMessage(msg *message.CLFMessage) error {

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
