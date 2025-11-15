package repository

import (
	"database/sql"
	"time"

	"github.com/Kost0/L4/internal/models"
)

func InsertEvent(db *sql.DB, event models.Event) error {
	query := `
	INSERT INTO events VALUES ($1, $2, $3, $4)
`

	_, err := db.Exec(query, event.EventID, event.UserID, event.Date, event.Event)
	if err != nil {
		return err
	}

	return nil
}

func UpdateEvent(db *sql.DB, event models.Event) error {
	query := `
	UPDATE events SET user_id = $1, date = $2, event = $3 WHERE event_id = $4
`

	_, err := db.Exec(query, event.UserID, event.Date, event.Event, event.EventID)
	if err != nil {
		return err
	}

	return nil
}

func DeleteEvent(db *sql.DB, id string) error {
	query := `
	DELETE FROM events WHERE event_id = $1
`

	_, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

func GetEventsByDate(db *sql.DB, date time.Time, days int, userID string) ([]*models.Event, error) {
	query := `
	SELECT * FROM events
	WHERE user_id = $1 && date BETWEEN $2 AND $3
`
	end := date.AddDate(0, 0, days)

	rows, err := db.Query(query, userID, date, end)
	if err != nil {
		return nil, err
	}

	var events []*models.Event

	for rows.Next() {
		event := models.Event{}
		err = rows.Scan(
			&event.EventID,
			&event.UserID,
			&event.Date,
			&event.Event,
		)
		if err != nil {
			return nil, err
		}

		events = append(events, &event)
	}

	return events, nil
}

func InsertToArchive(db *sql.DB, event models.Event) error {
	query := `
	INSERT INTO archive VALUES ($1, $2, $3, $4)
`

	_, err := db.Exec(query, event.EventID, event.UserID, event.Date, event.Event)
	if err != nil {
		return err
	}

	return nil
}
