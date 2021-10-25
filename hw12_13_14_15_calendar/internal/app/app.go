package app

import (
	"context"

	"github.com/DalerBakhriev/otus_go_homeworks/hw12_13_14_15_calendar/internal/store"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type App struct { // TODO
	logger *zap.Logger
	store  store.Store
}

func New(logger *zap.Logger, store store.Store) *App {
	return &App{logger: logger, store: store}
}

func NewDB(databaseURL string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
