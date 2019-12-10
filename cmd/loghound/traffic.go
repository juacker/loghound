package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

func main() {

	logfile := flag.String("l", "/tmp/access.log", "file to store logs")

	flag.Parse()

	f, err := os.Create(*logfile)
	check(err)

	defer f.Close()

	for {
		_, err = f.WriteString(getLog())
		if err != nil {
			log.Fatalf("fail writing file")
		}
		time.Sleep(10 * time.Duration(rand.Intn(10)) * time.Millisecond)
		f.Sync()
	}

}

var (
	users = []string{
		"Alice",
		"Bob",
		"Carol",
		"Carlos",
		"Charly",
		"Dan",
		"Erin",
		"Faythe",
	}

	methods = []string{
		"GET",
		"PUT",
		"POST",
		"DELTE",
	}

	paths = []string{
		"/admin/one",
		"/admin/two",
		"/admin/three",
		"/users/one",
		"/users/two",
		"/users/three",
		"/customers/one",
		"/customers/two",
		"/customers/three",
		"/ips/one",
		"/ips/two",
		"/ips/three",
		"/example/one",
		"/example/two",
		"/example/three",
		"/log/one",
		"/log/two",
		"/log/three",
		"/news/one",
		"/news/two",
		"/news/three",
		"/home/one",
		"/home/two",
		"/home/three",
	}

	status = []string{
		"200",
		"202",
		"300",
		"301",
		"400",
		"404",
		"422",
		"500",
	}
)

func getLog() string {
	now := time.Now().Format("02/Jan/2006:15:04:05 -0700")
	ip := fmt.Sprintf(
		"%d.%d.%d.%d",
		rand.Intn(256), rand.Intn(256), rand.Intn(256), rand.Intn(256),
	)

	user := users[rand.Intn(len(users)-1)]
	method := methods[rand.Intn(len(methods)-1)]
	path := paths[rand.Intn(len(paths)-1)]
	status := status[rand.Intn(len(status)-1)]
	bytes := rand.Intn(10000)

	return fmt.Sprintf("%s - %s [%s] \"%s %s HTTP/1.0\" %s %d\n", ip, user, now, method, path, status, bytes)

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
