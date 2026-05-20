package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/EduardoBacarin/golog-aggregator/src/server"
	"github.com/EduardoBacarin/golog-aggregator/src/writer"
)

func main() {
	logChannel := make(chan string, 100)
	fileChannel := make(chan string, 100)
	connType := "tcp"
	var wg sync.WaitGroup

	go writer.FileWriterWorker(fileChannel)
	if connType == "tcp" {
		go server.StartTCP(5000, logChannel)
	} else {
		go server.StartUDP(5000, logChannel)
	}

	fmt.Println("Ready to work!")

	numWorkers := 10
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			LogWorker(i, logChannel, fileChannel)
		}(i)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	fmt.Println("\n Closing. Waiting for pending processes")
	close(logChannel)

	wg.Wait()

	close(fileChannel)
}

func LogWorker(id int, logChannel <-chan string, fileChannel chan<- string) {
	for msg := range logChannel {
		processedMsg := fmt.Sprintf("[W-%d] %s", id, msg)
		fileChannel <- processedMsg
	}
}
