package postgres_storage_test

import (
	"database/sql"
	"fmt"
	"strings"
	"testing"
	"vox-server/internal/models"
	"vox-server/internal/storage/postgres_storage"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

// TODO : change initialization database_url (connect to config)
const database_url string = "host=localhost user=timno password=pass dbname=gitserver_test sslmode=disable"

func MakeTestDB(t *testing.T) (*sql.DB, func(...string)) {
	t.Helper()

	db, err := sql.Open("postgres", database_url)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	if err := db.Ping(); err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}

	return db, func(tables ...string) {
		if len(tables) > 0 {
			_, err := db.Exec(fmt.Sprintf("TRUNCATE %s CASCADE", strings.Join(tables, ", ")))
			if err != nil {
				t.Logf("Failed to truncate tables: %v", err)
			}
		}

		if err := db.Close(); err != nil {
			t.Logf("Failed to close database connection: %v", err)
		}
	}
}

func TestUserRepository_Create(t *testing.T) {
	db, cleanup := MakeTestDB(t)
	defer cleanup("users")

	storage := postgres_storage.NewDBStorage(db)
	repo := storage.Users()

	user := &models.User{
		Login:    "testuser",
		Username: "TestUser",
		Email:    "test@example.com",
		Password: "password",
	}

	err := repo.Create(user)

	assert.NoError(t, err, "Create should not return an error")
	assert.NotNil(t, user)

	foundUser, err := repo.FindByLogin("testuser")
	if err != nil {
		t.Fatalf("Failed to find user: %v", err)
	}

	if foundUser.Login != "testuser" {
		t.Errorf("Expected login 'testuser', got '%s'", foundUser.Login)
	}

}

func TestUserRepository_Count(t *testing.T) {
	db, cleanup := MakeTestDB(t)
	defer cleanup("users")

	storage := postgres_storage.NewDBStorage(db)
	repo := storage.Users()

	assert.Equal(t, 0, repo.Count(), "Expected 0 users in the database")

	user := &models.User{
		Login:    "testuser",
		Username: "TestUser",
		Email:    "test@example.com",
		Password: "password",
	}
	err := repo.Create(user)
	assert.NoError(t, err, "Create should not return an error")

	assert.Equal(t, 1, repo.Count(), "Expected 1 user in the database")
}

func TestUserRepository_IsEmpty(t *testing.T) {
	db, cleanup := MakeTestDB(t)
	defer cleanup("users")

	storage := postgres_storage.NewDBStorage(db)
	repo := storage.Users()

	assert.True(t, repo.IsEmpty(), "Expected database to be empty")

	user := &models.User{
		Login:    "testuser",
		Username: "TestUser",
		Email:    "test@example.com",
		Password: "password",
	}
	err := repo.Create(user)
	assert.NoError(t, err, "Create should not return an error")

	assert.False(t, repo.IsEmpty(), "Expected database to not be empty")
}

func TestUserRepository_FindByEmail(t *testing.T) {
	db, cleanup := MakeTestDB(t)
	defer cleanup("users")

	storage := postgres_storage.NewDBStorage(db)
	repo := storage.Users()

	user := &models.User{
		Login:    "testuser",
		Username: "TestUser",
		Email:    "test@example.com",
		Password: "password",
	}
	err := repo.Create(user)
	assert.NoError(t, err, "Create should not return an error")

	foundUser, err := repo.FindByEmail("test@example.com")
	assert.NoError(t, err, "FindByEmail should not return an error")
	assert.Equal(t, user.Login, foundUser.Login, "Expected login to match")
	assert.Equal(t, user.Email, foundUser.Email, "Expected email to match")
}

func TestUserRepository_DeleteByLogin(t *testing.T) {
	db, cleanup := MakeTestDB(t)
	defer cleanup("users")

	storage := postgres_storage.NewDBStorage(db)
	repo := storage.Users()

	user := &models.User{
		Login:    "testuser",
		Username: "TestUser",
		Email:    "test@example.com",
		Password: "password",
	}
	err := repo.Create(user)
	assert.NoError(t, err, "Create should not return an error")

	err = repo.DeleteByLogin("testuser")
	assert.NoError(t, err, "DeleteByLogin should not return an error")

	_, err = repo.FindByLogin("testuser")
	assert.Error(t, err, "Expected error when finding deleted user")
}

func TestUserRepository_DeleteByEmail(t *testing.T) {
	db, cleanup := MakeTestDB(t)
	defer cleanup("users")

	storage := postgres_storage.NewDBStorage(db)
	repo := storage.Users()

	user := &models.User{
		Login:    "testuser",
		Username: "TestUser",
		Email:    "test@example.com",
		Password: "password",
	}
	err := repo.Create(user)
	assert.NoError(t, err, "Create should not return an error")

	err = repo.DeleteByEmail("test@example.com")
	assert.NoError(t, err, "DeleteByEmail should not return an error")

	_, err = repo.FindByEmail("test@example.com")
	assert.Error(t, err, "Expected error when finding deleted user")
}

// TODO : test update-method when it is implemented
