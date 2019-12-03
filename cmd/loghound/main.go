package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/juacker/loghound/internal/alerts"
	"github.com/juacker/loghound/internal/broker"
	"github.com/juacker/loghound/internal/filemon"
	"github.com/juacker/loghound/internal/stats"
)

func main() {

	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool)

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()

	// goroutines control channel
	ctl := make(chan bool)

	var wg sync.WaitGroup

	wg.Add(4)

	go broker.Run(&wg, ctl)
	go filemon.Run(&wg, ctl)
	go stats.Run(&wg, ctl)

	go alerts.Run(&wg, ctl, "requests.total", "mean", 20, 1)

	<-done
	log.Println("main: Signal received, stopping goroutines")

	// stopping goroutines
	for i := 0; i < 4; i++ {
		ctl <- true
	}

	// waiting until they finish
	wg.Wait()
	log.Println("main: All goroutines stopped, exiting")
}
