package server_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"vox-server/internal/server"

	"github.com/stretchr/testify/assert"
)

// TODO : after TestServerWithDB_* testdb should be cleared
func /*Test*/ ServerWithDB_HandleUsersCreate(t *testing.T) {
	s, err := server.StartServer(true)
	if err != nil {
		log.Fatal(err)
	}

	testCases := []struct {
		name         string
		payload      any
		expectedCode int
	}{
		{
			name: "valid",
			payload: map[string]string{
				"login":    "login1",
				"username": "username",
				"email":    "user@example.org",
				"password": "password",
			},
			expectedCode: http.StatusCreated,
		},
		{
			name:         "invalid payload",
			payload:      "invalid",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid params[email]",
			payload: map[string]string{
				"login":    "login2",
				"username": "username",
				"email":    "user.",
				"password": "password",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "invalid params[empty password]",
			payload: map[string]string{
				"login":    "login3",
				"username": "username",
				"email":    "user@yandex.ru",
				"password": "",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)
			req, _ := http.NewRequest(http.MethodPost, "/users", b)

			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

// TODO : func TestServerWithDB_HandleSessionsCreate(t *testing.T)
// TODO : another tests
