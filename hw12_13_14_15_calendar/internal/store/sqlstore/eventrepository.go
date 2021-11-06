package sqlstore

import (
	"time"

	"github.com/DalerBakhriev/otus_go_homeworks/hw12_13_14_15_calendar/internal/model"
)

type EventRepository struct {
	store *Store
}

func (r *EventRepository) Create(e *model.Event) error {
	return r.store.db.QueryRow(
		`INSERT INTO events (title,
			                 creation_time,
							 start_time,
							 end_time,
							 description,
							 time_before_notification)
		 VALUES (:title, :creation_time, :start_time, :end_time, :description, :time_before_notification)
		 RETURNING id`,
		e.Title,
		e.CreationTime,
		e.StartTime,
		e.EndTime,
		e.Description,
		e.TimeBeforeNotification,
	).Scan(&e.ID)
}

func (r *EventRepository) Update(id int, e *model.Event) error {
	_, err := r.store.db.Exec(
		`UPDATE events
		 SET title = $1,
		     creation_time = $2,
			 start_time = $3,
			 end_time = $4,
			 description = $5,
			 time_before_notification = $6
		 WHERE id = $7`,
		e.Title,
		e.CreationTime,
		e.StartTime,
		e.EndTime,
		e.Description,
		e.TimeBeforeNotification,
		e.ID,
	)

	return err
}

func (r *EventRepository) Delete(id int) error {
	_, err := r.store.db.Exec(
		`DELETE
		 FROM events
		 WHERE id = $1`,
		id,
	)

	return err
}

func (r *EventRepository) ListByDay(date time.Time) ([]*model.Event, error) {
	events := make([]*model.Event, 0)
	err := r.store.db.Select(
		&events,
		`SELECT title,
		        creation_time,
				start_time,
				end_time,
				description,
				time_before_notification
		 FROM events
		 WHERE start_time = $1`,
		date,
	)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (r *EventRepository) ListByWeek(date time.Time) ([]*model.Event, error) {
	events := make([]*model.Event, 0)
	err := r.store.db.Select(
		&events,
		`SELECT title,
		        creation_time,
				start_time,
				end_time,
				description,
				time_before_notification
		 FROM events
		 WHERE start_time >= $1 and start_time <= start_time + make_interval(days => 7)`,
		date,
	)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (r *EventRepository) ListByMonth(date time.Time) ([]*model.Event, error) {
	events := make([]*model.Event, 0)
	err := r.store.db.Select(
		&events,
		`SELECT title,
		        creation_time,
				start_time,
				end_time,
				description,
				time_before_notification
		 FROM events
		 WHERE start_time >= $1 and start_time <= start_time + make_interval(months => 1)`,
		date,
	)
	if err != nil {
		return nil, err
	}

	return events, nil
}
