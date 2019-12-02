package alerts

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/juacker/loghound/internal/broker"
	"github.com/juacker/loghound/internal/message"
)

type metricMonitor struct {
	ctl             chan bool
	wg              *sync.WaitGroup
	broker          *broker.Connection
	cache           *cache
	metric          string
	operation       string
	threshold       float64
	interval        int64
	currentSeverity message.Severity
}

func (a *metricMonitor) loop() {
	log.Println("alerts: initializing alerts monitoring")

	ticker := time.NewTicker(5 * time.Second)

LOOP:
	for {
		select {
		case payload := <-a.broker.Receive():
			log.Println("alerts: new message received", payload)
			err := a.processMessage(payload)
			if err != nil {
				log.Println("alerts: failed processing message: ", err)
			}
		case <-ticker.C:
			log.Println("alerts: checking alerts")
			err := a.checkAlert()
			if err != nil {
				log.Println("alerts: failed cheking alerts: ", err)
			}
		case <-a.ctl:
			log.Println("alerts: ctl signal received, exiting")
			break LOOP
		}
	}

	a.wg.Done()
}

func (a *metricMonitor) processMessage(payload []byte) error {
	var msg message.StatMessage
	err := json.Unmarshal(payload, &msg)
	if err != nil {
		return err
	}

	// check message is the expected
	if !msg.IsValid() {
		return fmt.Errorf("invalid message")
	}

	return a.processStatMessage(&msg)
}

func (a *metricMonitor) processStatMessage(msg *message.StatMessage) error {

	a.cache.Add(
		Datapoint{
			Timestamp: time.Now().Unix(),
			Metric:    a.metric,
			Value:     msg.Stats[a.metric],
		},
	)

	return nil
}

func (a *metricMonitor) checkAlert() error {
	thresholdRaised := a.cache.Mean() > a.threshold
	text := fmt.Sprintf("High traffic generated an alert - hits = {%.2f}, triggered at {%v}", a.cache.Mean(), time.Now())

	var msg *message.AlertMessage
	if thresholdRaised && a.currentSeverity == message.SeverityCanceled {
		log.Println("alerts: new alert detected for metric ", a.metric)
		msg = message.NewAlertMessage(a.metric, text, message.SeverityMax)
	} else if !thresholdRaised && a.currentSeverity == message.SeverityMax {
		log.Println("alerts: cancelling alert for metric ", a.metric)
		msg = message.NewAlertMessage(a.metric, text, message.SeverityCanceled)
	} else {
		log.Println("alerts: nothing to do for alert ", a.metric)
		return nil
	}

	return a.broker.Send(broker.TopicAlert, msg)
}

// Run starts alerts
func Run(wg *sync.WaitGroup, ctl chan bool, metric, operation string, interval int64, threshold float64) {
	conn, err := broker.NewConnection(broker.TopicStat)
	if err != nil {
		log.Fatal("alerts: failed opening broker connection ", err)
	}

	p := &metricMonitor{
		ctl:       ctl,
		wg:        wg,
		broker:    conn,
		metric:    metric,
		operation: operation,
		interval:  interval,
		threshold: threshold,
		cache: &cache{
			elements: make(Elements, 0),
		},
	}

	p.loop()
}
