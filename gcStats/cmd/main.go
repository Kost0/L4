package main

import (
	"errors"
	"log"
	"net/http"
	_ "net/http/pprof"
	"sync"

	"github.com/Kost0/L4/internal/analytics"
	"github.com/Kost0/L4/internal/handlers"
	"github.com/Kost0/L4/internal/someWork"
)

func main() {
	listenPort := ":8080"
	pprofPort := ":8081"

	go analytics.Worker()

	http.HandleFunc("/stats", handlers.GetStats)

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := http.ListenAndServe(listenPort, nil); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Error starting HTTP server: %v\n", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := http.ListenAndServe(pprofPort, nil); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Error starting pprof HTTP server: %v\n", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		someWork.StartWork()
	}()

	wg.Wait()
}
