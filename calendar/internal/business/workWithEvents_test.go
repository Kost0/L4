package business

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/Kost0/L4/internal/models"
	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"
)

var (
	userID = uuid.New()
	date   = time.Now()
	repo   = EventRepository{DB: nil}
)

func TestNewEvent_Success(t *testing.T) {
	event := &models.EventDTO{
		UserID: userID,
		Date:   date,
		Event:  "test event",
	}

	buf, err := json.Marshal(event)
	assert.NoError(t, err)

	defer func() {
		if r := recover(); r == nil {
			t.Error("Execution did not reach the point of insertion into db.")
		}
	}()

	repo.NewEvent(buf)
}

func TestNewEvent_WrongStruct(t *testing.T) {
	event := struct {
		Number int
	}{
		Number: 1,
	}

	buf, err := json.Marshal(event)
	assert.NoError(t, err)

	_, err = repo.NewEvent(buf)
	assert.Error(t, err)
}

func TestUpdateEvent_Success(t *testing.T) {
	models.Events = []*models.Event{{
		EventID: uuid.New(),
		UserID:  userID,
		Date:    date,
		Event:   "test event",
	}}

	event := &models.EventDTO{
		UserID: userID,
		Date:   date.Add(24 * time.Hour),
		Event:  "test event",
	}

	buf, err := json.Marshal(event)
	assert.NoError(t, err)

	defer func() {
		if r := recover(); r == nil {
			t.Error("Execution did not reach the point of updating into db.")
		}
	}()

	repo.UpdateEvent(buf, models.Events[0].EventID.String())
}

func TestUpdateEvent_NothingToUpdate(t *testing.T) {
	event := &models.Event{
		EventID: uuid.New(),
		UserID:  userID,
		Date:    date,
		Event:   "test event",
	}

	buf, err := json.Marshal(event)
	assert.NoError(t, err)

	updatedEvent, err := repo.UpdateEvent(buf, "1")
	assert.Error(t, err)
	assert.Nil(t, updatedEvent)
}

func TestUpdateEvent_WrongStruct(t *testing.T) {
	models.Events = []*models.Event{{
		UserID: userID,
		Date:   date,
		Event:  "test event",
	}}

	event := struct {
		Number int
	}{
		Number: 1,
	}

	buf, err := json.Marshal(event)
	assert.NoError(t, err)

	updatedEvent, err := repo.UpdateEvent(buf, "1")
	assert.Error(t, err)
	assert.Nil(t, updatedEvent)
}

func TestDeleteEvent_Success(t *testing.T) {
	models.Events = []*models.Event{{
		EventID: uuid.New(),
		UserID:  userID,
		Date:    date,
		Event:   "test event",
	}}

	defer func() {
		if r := recover(); r == nil {
			t.Error("Execution did not reach the point of deleting into db.")
		}
	}()

	repo.DeleteEvent(models.Events[0].EventID.String())
}

func TestDeleteEvent_NothingToDelete(t *testing.T) {
	err := repo.DeleteEvent("1")
	assert.Error(t, err)
}
