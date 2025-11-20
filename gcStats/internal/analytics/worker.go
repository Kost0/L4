package analytics

import (
	"runtime"
	"time"

	"github.com/Kost0/L4/internal/models"
)

func Worker() {
	timer := time.NewTimer(15 * time.Second)

	for {
		<-timer.C

		stats := getAnalytics()

		models.AllStats = append(models.AllStats, stats)

		timer.Reset(15 * time.Second)
	}
}

func getAnalytics() models.Stats {
	stats := runtime.MemStats{}

	runtime.ReadMemStats(&stats)

	malloc := stats.Mallocs
	numGc := stats.NumGC
	heapAlloc := stats.HeapAlloc
	lastGC := stats.LastGC

	newStats := models.Stats{
		Malloc:    malloc,
		NumGC:     numGc,
		HeapAlloc: heapAlloc,
		LastGC:    lastGC,
		Time:      time.Now(),
	}

	return newStats
}
