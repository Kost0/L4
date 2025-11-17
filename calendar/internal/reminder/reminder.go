package reminder

import (
	"container/heap"
	"context"
	"log"
	"math"
	"sync"
	"time"

	"github.com/Kost0/L4/internal/models"
)

func Worker(ctx context.Context, ch chan *models.Event, wg *sync.WaitGroup) {
	defer wg.Done()

	eventHeap := CreateHeap()

	var timer *time.Timer
	var timerChan <-chan time.Time

	for {
		select {
		case event := <-ch:
			heap.Push(eventHeap, event)

			earliestTask := (*eventHeap)[0]

			if timer == nil {
				timer = time.NewTimer(time.Until(earliestTask.RemindAt))
			} else {
				if !timer.Stop() {
					<-timer.C
				}
				timer.Reset(time.Until(earliestTask.RemindAt))
			}

			timerChan = timer.C
		case <-timerChan:
			event := heap.Pop(eventHeap).(*models.Event)

			log.Printf("Напоминание, до события %s осталось %v минут(а)", event.Event, math.Ceil(time.Until(event.Date).Minutes()))

			if eventHeap.Len() > 0 {
				nextEvent := (*eventHeap)[0]
				timer.Reset(time.Until(nextEvent.RemindAt))
				timerChan = timer.C
			} else {
				timerChan = nil
			}
		case <-ctx.Done():
			log.Println("context done, worker stopped")
			return
		}
	}

}
