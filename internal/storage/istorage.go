package storage

type Storage interface {
	Users() UserRepository
}
