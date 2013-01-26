package producer

import (
	"log"
)

var logger *log.Logger

func Init(l *log.Logger) {
	logger = l
	logger.Println("[Consumer] Started")
}

func Tick(t int) {

}
