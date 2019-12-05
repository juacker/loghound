package main

import (
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

	// print logs to file
	f, err := os.OpenFile("loghound.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening log file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	// goroutines control channel
	ctl := make(chan bool)

	var wg sync.WaitGroup

	wg.Add(4)

	go broker.Run(&wg, ctl)
	go filemon.Run(&wg, ctl)
	go stats.Run(&wg, ctl)

	go alerts.Run(&wg, ctl, "requests.total", "mean", 20, 1)

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
