package sqlstore

import (
	"context"

	"github.com/DalerBakhriev/otus_go_homeworks/hw12_13_14_15_calendar/internal/store"
	"github.com/jmoiron/sqlx"
)

type Store struct {
	db              *sqlx.DB
	eventRepository *EventRepository
}

func New(db *sqlx.DB) *Store {
	return &Store{db: db}
}

func (s *Store) Event() store.EventRepository {
	if s.eventRepository != nil {
		return s.eventRepository
	}
	return &EventRepository{store: s}
}

func (s *Store) Connect(ctx context.Context) error {
	db, err := sqlx.Connect("postgres", "user=foo dbname=bar sslmode=disable")
	if err != nil {
		return err
	}
	s.db = db
	return nil
}

func (s *Store) Close(ctx context.Context) error {
	// TODO
	return nil
}
