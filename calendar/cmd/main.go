package main

import (
	"context"
	"errors"
	"os/signal"
	"sync"
	"syscall"

	"log"
	"net/http"
	"os"

	"github.com/Kost0/L4/internal/business"
	"github.com/Kost0/L4/internal/collector"
	"github.com/Kost0/L4/internal/handlers"
	"github.com/Kost0/L4/internal/middleware"
	"github.com/Kost0/L4/internal/models"
	"github.com/Kost0/L4/internal/reminder"
	"github.com/Kost0/L4/internal/repository"

	"github.com/go-chi/chi/v5"
)

func main() {
	db, err := repository.ConnectDB()
	if err != nil {
		panic(err)
	}

	err = repository.RunMigrations(db, "events_db")
	if err != nil {
		panic(err)
	}

	eventRepository := business.EventRepository{
		DB: db,
	}

	ch := make(chan *models.Event)

	logCh := make(chan string)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	col := collector.Collector{
		DB: db,
	}

	logger := middleware.Logger{
		Ch: logCh,
	}

	wg := &sync.WaitGroup{}

	wg.Add(3)

	go logger.StartLogger(ctx, wg)

	go col.StartCollector(ctx, wg)

	go reminder.Worker(ctx, ch, wg)

	handler := handlers.Handler{
		Ch:    ch,
		LogCh: logCh,
		Rep:   &eventRepository,
	}

	r := chi.NewRouter()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	r.Use(logger.Middleware)

	r.Post("/delete_event", handler.DeleteEvent)

	r.Post("/update_event", handler.UpdateEvent)

	r.Post("/create_event", handler.CreateEvent)

	r.Get("/events_for_day", handler.GetEventForDay)

	r.Get("/events_for_week", handler.GetEventForWeek)

	r.Get("/events_for_month", handler.GetEventForMonth)

	log.Println("Starting HTTP server...")

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Error starting HTTP server: %v\n", err)
		}
	}()

	wg.Wait()
}
