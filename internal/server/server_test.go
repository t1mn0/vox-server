package server_test

import (
	"git-server/internal/server"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServer_HandleUsersCreate(t *testing.T) {
	// в ближайшем будущем переделать
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/users", nil)
	cfg := server.NewConfig()
	cfg.DatabaseURL = "postgres://timno:pass@localhost/git-server_test?sslmode=disable"

	// TODO : change to local storage (NewInMemoryServer)
	server, err := server.NewServerWithDB(cfg)
	if err != nil {
		log.Fatal(err)
	}

	server.ServeHTTP(rec, req)
	assert.Equal(t, rec.Code, http.StatusOK)
}
