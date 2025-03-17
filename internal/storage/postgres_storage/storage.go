package postgres_storage

import (
	"database/sql"
	"vox-server/internal/storage"

	_ "github.com/lib/pq"
)

type DBStorage struct {
	db *sql.DB
}

func NewDBStorage(db *sql.DB) *DBStorage {
	return &DBStorage{
		db: db,
	}
}

func (storage *DBStorage) Users() storage.UserRepository {
	return UserRepository{storage: storage}
}
