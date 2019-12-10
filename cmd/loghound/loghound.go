package main

import (
	"flag"
	"log"
	"os"
	"sync"

	"github.com/juacker/loghound/internal/alerts"
	"github.com/juacker/loghound/internal/broker"
	"github.com/juacker/loghound/internal/console"
	"github.com/juacker/loghound/internal/filemon"
	"github.com/juacker/loghound/internal/stats"
)

func main() {

	// parse command line arguments
	logfile := flag.String("l", "/tmp/access.log", "common log format file to monitor")
	threshold := flag.Int("t", 10, "alarm threshold (req/seq)")
	alarmInterval := flag.Int64("a", 120, "interval to consider for alarm threshold (s)")
	statsInterval := flag.Int64("s", 2, "stats interval generation (s)")

	flag.Parse()

	// print logs to file
	f, err := os.OpenFile("loghound.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening log file: %v", err)
	}
	defer f.Close()

	log.Println("writing logs to ./loghound.log file")

	log.SetOutput(f)

	// goroutines control channel
	ctl := make(chan bool)

	var wg sync.WaitGroup

	wg.Add(4)

	go broker.Run(&wg, ctl)
	go filemon.Run(&wg, ctl, *logfile)
	go stats.Run(&wg, ctl, *statsInterval)

	go alerts.Run(&wg, ctl, "requests.total", "mean", *alarmInterval, *threshold)

	console.Run()

	log.Println("main: stopping goroutines")

	// stopping goroutines
	for i := 0; i < 4; i++ {
		ctl <- true
	}

	// waiting until they finish
	wg.Wait()
	log.Println("main: All goroutines stopped, exiting")
}
