package business

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/Kost0/L4/internal/models"
	"github.com/Kost0/L4/internal/repository"
	"github.com/google/uuid"
)

type EventRepository struct {
	DB *sql.DB
}

func (e *EventRepository) NewEvent(buf []byte) (*models.Event, error) {
	eventDTO := &models.EventDTO{}

	decoder := json.NewDecoder(bytes.NewReader(buf))

	decoder.DisallowUnknownFields()

	err := decoder.Decode(eventDTO)
	if err != nil {
		return nil, err
	}

	event := &models.Event{
		EventID:  uuid.New(),
		UserID:   eventDTO.UserID,
		Date:     eventDTO.Date,
		RemindAt: eventDTO.Date.Add(-2 * time.Hour),
		Event:    eventDTO.Event,
	}

	if event.Date.Before(time.Now()) {
		return nil, errors.New("date is in the past")
	}

	models.Events = append(models.Events, event)

	log.Printf("Appended event %v", event)

	err = repository.InsertEvent(e.DB, *event)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (e *EventRepository) UpdateEvent(buf []byte, id string) (*models.Event, error) {
	eventDTO := &models.EventDTO{}

	decoder := json.NewDecoder(bytes.NewReader(buf))

	decoder.DisallowUnknownFields()

	err := decoder.Decode(eventDTO)
	if err != nil {
		return nil, err
	}

	eventUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	event := &models.Event{
		EventID:  eventUUID,
		UserID:   eventDTO.UserID,
		Date:     eventDTO.Date,
		RemindAt: eventDTO.Date.Add(-2 * time.Hour),
		Event:    eventDTO.Event,
	}

	hit := false

	for _, ev := range models.Events {
		log.Printf("events: %v", ev)
	}

	for i, ev := range models.Events {
		if ev.EventID == event.EventID {
			models.Events[i] = event
			hit = true
			break
		}
	}

	if hit {
		err = repository.UpdateEvent(e.DB, *event)
		if err != nil {
			return nil, err
		}

		return event, nil
	}

	return nil, errors.New("event not found")
}

func (e *EventRepository) DeleteEvent(id string) error {
	hit := false

	for i, ev := range models.Events {
		if ev.EventID.String() == id {
			models.Events = append(models.Events[:i], models.Events[i+1:]...)
			hit = true
		}
	}

	if hit {
		err := repository.DeleteEvent(e.DB, id)
		if err != nil {
			return err
		}

		return nil
	}

	return errors.New("event not found")
}

func (e *EventRepository) FindEventsForTime(date time.Time, userID string, days int) ([]*models.Event, error) {
	day := date.Truncate(24 * time.Hour)
	return repository.GetEventsByDate(e.DB, day, days, userID)
}

func (e *EventRepository) WakeUp() {
	var err error

	models.Events, err = repository.GetAllEvents(e.DB)
	if err != nil {
		log.Println(err)
	}

	log.Printf("%d events will be waked", len(models.Events))
}
