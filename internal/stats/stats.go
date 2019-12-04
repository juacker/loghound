package stats

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/juacker/loghound/internal/broker"
	"github.com/juacker/loghound/internal/message"
)

type statsMonitor struct {
	ctl    chan bool
	wg     *sync.WaitGroup
	broker *broker.Connection
	cache  *cache
}

func (s *statsMonitor) loop() {
	log.Println("stats: initializing stats monitoring")

	ticker := time.NewTicker(10 * time.Second)

LOOP:
	for {
		select {
		case payload := <-s.broker.Receive():
			log.Println("stats: new message received", payload)
			err := s.processMessage(payload)
			if err != nil {
				log.Println("stats: failed processing message: ", err)
			}
		case <-ticker.C:
			log.Println("stats: sending stats")
			err := s.sendStats()
			if err != nil {
				log.Println("stats: failed sending stats: ", err)
			}
		case <-s.ctl:
			log.Println("stats: ctl signal received, exiting")
			break LOOP
		}
	}

	s.wg.Done()
}

func (s *statsMonitor) processMessage(payload []byte) error {
	var msg message.CLFMessage
	err := json.Unmarshal(payload, &msg)
	if err != nil {
		return err
	}

	// check message is the expected
	if !msg.IsValid() {
		return fmt.Errorf("invalid message")
	}

	return s.processCLFMessage(&msg)
}

func (s *statsMonitor) processCLFMessage(msg *message.CLFMessage) error {

	// Process message fields
	// root path
	paths := strings.Split(msg.Request.Path, "/")
	var rootPath string
	if len(paths) == 0 {
		return fmt.Errorf("invalid path detected")
	} else if len(paths) == 1 {
		rootPath = "/"
	} else {
		rootPath = "/" + paths[1]
	}

	// status
	status := string(msg.Status)

	// create some counter metrics using message fields

	// metric: requests.total
	s.cache.Increment("requests.total", 1)

	// metric: bytes.total
	s.cache.Increment("bytes.total", msg.Bytes)

	// metric: path.<path>.requests
	s.cache.Increment("path."+rootPath+".requests", 1)

	// metric: path.<path>.bytes
	s.cache.Increment("path."+rootPath+".bytes", msg.Bytes)

	// metric: status.<status>.requests
	s.cache.Increment("status."+status+".requests", 1)

	// metric: status.<status>.bytes
	s.cache.Increment("status."+status+".bytes", msg.Bytes)

	// metric: method.<method>.requests
	s.cache.Increment("method."+msg.Request.Method+".requests", 1)

	// metric: method.<method>.bytes
	s.cache.Increment("method."+msg.Request.Method+".bytes", msg.Bytes)

	return nil
}

func (s *statsMonitor) sendStats() error {

	stats, begin, end := s.cache.Stats()
	return s.broker.Send(broker.TopicStat, message.NewStatMessage(stats, begin, end))
}

// Run starts stats
func Run(wg *sync.WaitGroup, ctl chan bool) {
	conn, err := broker.NewConnection(broker.TopicData)
	if err != nil {
		log.Fatal("stats: failed opening broker connection ", err)
	}

	stats := &statsMonitor{
		ctl:    ctl,
		wg:     wg,
		broker: conn,
		cache: &cache{
			metrics: make(map[string]int),
			reset:   time.Now().Unix(),
		},
	}

	stats.loop()
}
