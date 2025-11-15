package middleware

import (
	"context"
	"log"
	"sync"
)

type Logger struct {
	Ch chan string
}

func (l *Logger) StartLogger(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case logStr := <-l.Ch:
			log.Println(logStr)
		case <-ctx.Done():
			return
		}
	}
}
