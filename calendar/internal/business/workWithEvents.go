package business

import (
	"database/sql"
	"encoding/json"
	"errors"
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

	err := json.Unmarshal(buf, eventDTO)
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

	models.Events = append(models.Events, event)

	err = repository.InsertEvent(e.DB, *event)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (e *EventRepository) UpdateEvent(buf []byte, id string) (*models.Event, error) {
	eventDTO := &models.EventDTO{}

	err := json.Unmarshal(buf, eventDTO)
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

	for i, ev := range models.Events {
		if ev.EventID == event.EventID {
			if ev.Event == event.Event || ev.Date == event.Date {
				models.Events[i] = event
				hit = true
				break
			}
		}
	}

	if hit {
		err = repository.UpdateEvent(e.DB, *event)
		if err != nil {
			return nil, err
		}
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
	}

	return errors.New("event not found")
}

func (e *EventRepository) FindEventsForTime(date time.Time, userID string, days int) ([]*models.Event, error) {
	return repository.GetEventsByDate(e.DB, date, days, userID)
}
