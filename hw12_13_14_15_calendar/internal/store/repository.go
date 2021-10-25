package store

import "github.com/DalerBakhriev/otus_go_homeworks/hw12_13_14_15_calendar/internal/model"

type EventRepository interface {
	Create(e *model.Event) error   // добавление события  вхранилище
	Update(e *model.Event) error   // изменение события в хранилище;
	Delete(id int) error           // удаление события из хранилища;
	List() ([]*model.Event, error) // листинг событий;
}
