package infrastructure

import (
	"log"
)

type Logger struct{}

func (l Logger) Log(args ...interface{}) {
	log.Println(args...)
}
