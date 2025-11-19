package main

import (
	"errors"
	"log"
	"net/http"
	"net/http/pprof"
	_ "net/http/pprof"
	"sync"

	"github.com/Kost0/L4/internal/handler"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	r.Post("/sort", handler.SortNums)

	r2 := chi.NewRouter()

	srv2 := &http.Server{
		Addr:    ":8081",
		Handler: r2,
	}

	r2.HandleFunc("/debug/pprof/*", pprof.Index)

	r2.HandleFunc("/debug/pprof/profile", pprof.Profile)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Error starting HTTP server: %v\n", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := srv2.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Error starting HTTP server: %v\n", err)
		}
	}()

	wg.Wait()
}
