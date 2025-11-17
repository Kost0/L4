package collector

import (
	"context"
	"database/sql"
	"log"
	"sync"
	"time"

	"github.com/Kost0/L4/internal/models"
	"github.com/Kost0/L4/internal/repository"
)

type Collector struct {
	DB *sql.DB
}

const interval = 2 * time.Minute

func (c *Collector) StartCollector(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	timer := time.NewTimer(interval)

	for {
		select {
		case <-timer.C:
			c.collect()

			timer.Reset(interval)
		case <-ctx.Done():
			log.Println("context done, collector stopped")
			return
		}
	}
}

func (c *Collector) collect() {
	log.Println("Collecting events...")

	old := len(models.Events)

	for i := len(models.Events) - 1; i >= 0; i-- {
		event := models.Events[i]
		if event.RemindAt.Before(time.Now()) {
			if i != len(models.Events)-1 {
				models.Events = append(models.Events[:i], models.Events[i+1:]...)
			} else {
				models.Events = models.Events[:i]
			}

			err := repository.DeleteEvent(c.DB, event.EventID.String())
			if err != nil {
				log.Println(err)
				continue
			}

			err = repository.InsertToArchive(c.DB, *event)
			if err != nil {
				log.Println(err)
			}
		}
	}

	log.Printf("%d evets deleted", old-len(models.Events))
}
