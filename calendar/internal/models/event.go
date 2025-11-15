package models

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	EventID  uuid.UUID `json:"event_id"`
	UserID   uuid.UUID `json:"user_id"`
	Date     time.Time `json:"date"`
	RemindAt time.Time `json:"-"`
	Event    string    `json:"event"`
}

type EventDTO struct {
	UserID uuid.UUID `json:"user_id"`
	Event  string    `json:"event"`
	Date   time.Time `json:"date"`
}

var Events = make([]*Event, 0)
