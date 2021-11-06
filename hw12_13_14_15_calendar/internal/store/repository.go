package store

import (
	"time"

	"github.com/DalerBakhriev/otus_go_homeworks/hw12_13_14_15_calendar/internal/model"
)

type EventRepository interface {
	Create(e *model.Event) error                        // добавление события  вхранилище
	Update(id int, e *model.Event) error                // изменение события в хранилище;
	Delete(id int) error                                // удаление события из хранилища;
	ListByDay(date time.Time) ([]*model.Event, error)   // листинг событий за день
	ListByWeek(date time.Time) ([]*model.Event, error)  // листинг событий за неделю
	ListByMonth(date time.Time) ([]*model.Event, error) // листинг событий за месяц
}
