package hub

import (
	"fmt"
	"log"
	"os"
	"time"
)

func Logger(message string, args ...any) {
	file, err := os.OpenFile("./log/chroma_hub.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Printf("Error opening log file (%s)", err)
	}

	t := time.Now()
	s := fmt.Sprintf("[%s]\t", t.Format("2006-01-02 15:04:05")) + fmt.Sprintf(message, args...) + "\n"
	file.Write([]byte(s))
}
