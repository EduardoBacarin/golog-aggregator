package writer

import (
	"fmt"
	"os"
	"time"
)

func FileWriterWorker(fileChan <-chan string) {
	file, err := os.OpenFile("log.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error on opening file: %v\n", err)
		return
	}
	defer file.Close()

	for msg := range fileChan {
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		entry := fmt.Sprintf("%s | %s\n", timestamp, msg)

		_, err := file.WriteString(entry)
		if err != nil {
			fmt.Printf("Error on writing file: %v\n", err)
		}
	}
}
