package store

type Store interface {
	Event() EventRepository
}
