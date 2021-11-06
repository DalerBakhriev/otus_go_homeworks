package memorystore

import (
	"sync"
	"time"

	"github.com/DalerBakhriev/otus_go_homeworks/hw12_13_14_15_calendar/internal/model"
)

type EventRepository struct {
	events map[int]*model.Event
	mu     sync.RWMutex
}

func (r *EventRepository) Create(e *model.Event) error {
	e.ID = len(r.events) + 1
	r.mu.Lock()
	r.events[e.ID] = e
	r.mu.Unlock()

	return nil
}

func (r *EventRepository) Update(id int, e *model.Event) error {
	e.ID = id
	r.mu.Lock()
	r.events[id] = e
	r.mu.Unlock()

	return nil
}

func (r *EventRepository) Delete(id int) error {
	r.mu.Lock()
	delete(r.events, id)
	r.mu.Unlock()

	return nil
}

func (r *EventRepository) ListByDay(date time.Time) ([]*model.Event, error) {
	events := make([]*model.Event, 0)
	for _, event := range r.events {
		if event.StartTime.Equal(date) {
			events = append(events, event)
		}
	}
	return events, nil
}

func (r *EventRepository) ListByWeek(date time.Time) ([]*model.Event, error) {
	events := make([]*model.Event, 0)
	for _, event := range r.events {
		dateWeekLater := date.AddDate(0, 0, 7)
		if (event.StartTime.After(date) || event.StartTime.Equal(date)) &&
			(event.StartTime.Before(dateWeekLater) || event.StartTime.Equal(dateWeekLater)) {
			events = append(events, event)
		}
	}
	return events, nil
}

func (r *EventRepository) ListByMonth(date time.Time) ([]*model.Event, error) {
	events := make([]*model.Event, 0)
	for _, event := range r.events {
		dateMonthLater := date.AddDate(0, 1, 0)
		if (event.StartTime.After(date) || event.StartTime.Equal(date)) &&
			(event.StartTime.Before(dateMonthLater) || event.StartTime.Equal(dateMonthLater)) {
			events = append(events, event)
		}
	}
	return events, nil
}
