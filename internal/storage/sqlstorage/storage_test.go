package sqlstorage_test

import (
	"os"
	"testing"
)

var databaseURL string

func TestMain(m *testing.M) {
	databaseURL = os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://timno:pass@localhost/git-server_test?sslmode=disable"
	}

	os.Exit(m.Run())
}
