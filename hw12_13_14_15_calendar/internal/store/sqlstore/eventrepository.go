package sqlstore

import (
	"github.com/DalerBakhriev/otus_go_homeworks/hw12_13_14_15_calendar/internal/model"
)

type EventRepository struct {
	store *Store
}

func (r *EventRepository) Create(e *model.Event) error {
	return r.store.db.QueryRow(
		`INSERT INTO events (title, creation_time, plan_time)
		 VALUES (:title, :creation_time, :plan_time)
		 RETURNING id`,
		e.Title,
		e.CreationTime,
		e.PlanTime,
	).Scan(&e.ID)
}

func (r *EventRepository) Update(e *model.Event) error {
	_, err := r.store.db.Exec(
		`UPDATE events
		 SET title = $1,
		     creation_time = $2,
			 plan_time = $3
		 WHERE id = $4`,
		e.Title,
		e.CreationTime,
		e.PlanTime,
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

func (r *EventRepository) List() ([]*model.Event, error) {
	events := make([]*model.Event, 0)
	err := r.store.db.Select(&events, "SELECT * FROM events")
	if err != nil {
		return nil, err
	}

	return events, nil
}
