package sqlstorage

import (
	"database/sql"
	"git-server/internal/storage"

	_ "github.com/lib/pq"
)

type DBStorage struct {
	db *sql.DB
}

func New(db *sql.DB) *DBStorage {
	return &DBStorage{
		db: db,
	}
}

func (st *DBStorage) User() storage.UserRepository {
	return UserRepo{storage: st}
}
