package memorystore

import (
	"sync"

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

func (r *EventRepository) Update(e *model.Event) error {
	r.mu.Lock()
	r.events[e.ID] = e
	r.mu.Unlock()

	return nil
}

func (r *EventRepository) Delete(id int) error {
	r.mu.Lock()
	delete(r.events, id)
	r.mu.Unlock()

	return nil
}

func (r *EventRepository) List() ([]*model.Event, error) {
	events := make([]*model.Event, 0, len(r.events))
	for _, event := range r.events {
		events = append(events, event)
	}
	return events, nil
}
