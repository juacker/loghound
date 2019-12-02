package filemon

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/juacker/loghound/internal/broker"
	"github.com/juacker/loghound/internal/message"
	"github.com/juacker/loghound/pkg/clf"
)

type fileMonitor struct {
	sync.Mutex
	ctl     chan bool
	wg      *sync.WaitGroup
	broker  *broker.Connection
	files   []string
	fd      map[string]*os.File
	watcher *fsnotify.Watcher
}

func (f *fileMonitor) loop() {
	log.Println("filemon: initializing file monitoring")

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	f.watcher = watcher
	defer f.watcher.Close()

	for _, filename := range f.files {
		file, err := os.Open(filename)
		if err != nil {
			log.Fatal("filemon: failed openning file: ", filename)
		}

		defer file.Close()
		f.fd[filename] = file

		// let's position the file descriptor at the end of the file
		// we want to process only new contents
		position, err := file.Seek(0, 2)
		if err != nil {
			log.Fatal("filemon: failed positioning at the end of file: ", filename)
		}
		log.Println("positioned at the end of file: ", filename, position)

		log.Println("filemon: adding file to monitoring list: ", filename)
		err = watcher.Add(filename)
		if err != nil {
			log.Fatal("filemon: could not add file ", filename, ", ", err)
		}
	}

LOOP:
	for {
		select {
		case event := <-watcher.Events:
			log.Println("filemon: new event received ", event, event.Name)
			if event.Op&fsnotify.Write == fsnotify.Write {
				f.processFileContents(event.Name)

			}
		case err := <-watcher.Errors:
			log.Println("filemon: error:", err)
		case <-f.ctl:
			log.Println("filemon: ctl signal received, exiting")
			break LOOP
		}
	}

	f.wg.Done()
}

func (f *fileMonitor) processFileContents(filename string) error {
	fd := f.fd[filename]
	if fd == nil {
		return fmt.Errorf("filemon: file descriptor not found for %s", filename)
	}

	// let's see if previous position is still valid
	// or file has been truncated
	previousPosition, err := fd.Seek(0, 1)
	if err != nil {
		return fmt.Errorf("filemon: fail seeking current position of file %s", filename)
	}
	eofPosition, err := fd.Seek(0, 2)
	if err != nil {
		return fmt.Errorf("filemon: fail seeking EOF position of file %s", filename)
	}

	if eofPosition < previousPosition {
		// file truncated let's move to the beginning of the file
		log.Println("filemon: file truncation detected: ", filename)
		fd.Seek(0, 0) //errcheck: nolint
	} else {
		// file not truncated let's move to the previous position
		fd.Seek(previousPosition-eofPosition, 1)
	}

	// Start reading from the file with a reader.
	// it starts where file descriptor was positioned previously
	reader := bufio.NewReader(fd)

	var line string
	for {
		line, err = reader.ReadString('\n')
		if err != nil {
			break
		}

		if len(line) > 0 {
			logEntry, err := clf.Parse(line)
			if err != nil {
				log.Println("filemon: failed parsing line for file ", filename, err)
				continue
			}

			err = f.broker.Send(broker.TopicData, message.NewCLFMessage(logEntry))
			if err != nil {
				log.Println("filemon: failed sending message to broker")
			}
		}
	}

	if err != io.EOF {
		log.Println(fmt.Sprintf("filemon: failed reading file %s: %v", filename, err))
		return err
	}

	return nil
}

// Run starts file monitor
func Run(wg *sync.WaitGroup, ctl chan bool) {
	conn, err := broker.NewConnection()
	if err != nil {
		log.Fatal("filemon: failed opening broker connection ", err)
	}

	filemon := &fileMonitor{
		ctl:    ctl,
		wg:     wg,
		files:  []string{"/tmp/access.log"},
		fd:     make(map[string]*os.File, 0),
		broker: conn,
	}

	filemon.loop()
}
