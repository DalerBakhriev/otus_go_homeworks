package memorystore

import (
	"github.com/DalerBakhriev/otus_go_homeworks/hw12_13_14_15_calendar/internal/model"
	"github.com/DalerBakhriev/otus_go_homeworks/hw12_13_14_15_calendar/internal/store"
)

type Store struct {
	eventRepository *EventRepository
}

func New() store.Store {
	return &Store{}
}

func (s *Store) Event() store.EventRepository {
	if s.eventRepository != nil {
		return s.eventRepository
	}
	s.eventRepository = &EventRepository{
		events: make(map[int]*model.Event),
	}
	return s.eventRepository
}
