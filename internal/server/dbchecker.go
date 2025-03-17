package server

import (
	"database/sql"
	"fmt"
	"os/exec"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var dbNames = []string{"gitserver", "gitserver_test"}

func isPsqlInstalled() bool {
	_, err := exec.LookPath("psql")
	return err == nil
}

func isMigrateInstalled() bool {
	_, err := exec.LookPath("migrate")
	return err == nil
}

func checkDBsExists(dbName string) (bool, error) {
	cmd := exec.Command("psql", "-lqt")
	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("failed to list databases: %v", err)
	}

	if strings.Contains(string(output), dbName) {
		return true, nil
	}
	return false, nil
}

func createDB(cfg *Config, dbName string) error {
	cmd := exec.Command("createdb", "-U", cfg.DB.User, "-h", cfg.DB.Host, "-p", cfg.DB.Port, dbName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create database %s: %v, output: %s", dbName, err, string(output))
	}
	return nil
}

func runMigrations(dbURL string) error {
	migrationsPath := "file://migrations"

	m, err := migrate.New(migrationsPath, dbURL)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %v", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			fmt.Println("No new migrations to apply.")
		} else {
			return fmt.Errorf("failed to apply migrations: %v", err)
		}
	} else {
		fmt.Println("Migrations applied successfully.")
	}

	return nil
}

func validateConfigDBData(cfg *Config) error {
	dbURL := cfg.DatabaseURL
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return fmt.Errorf("connection error to the main database: %w", err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		return fmt.Errorf("couldn't connect to the main database: %w", err)
	}

	testDBURL := cfg.TestDatabaseURL
	testDB, err := sql.Open("postgres", testDBURL)
	if err != nil {
		return fmt.Errorf("error connecting to the test database: %w", err)
	}
	defer testDB.Close()
	if err := testDB.Ping(); err != nil {
		return fmt.Errorf("couldn't connect to the test database: %w", err)
	}

	return nil
}

func ConfigurationDBs(cfg *Config) error {
	if !isPsqlInstalled() {
		return fmt.Errorf("psql is not installed. Please install PostgreSQL to continue")
	}

	if !isMigrateInstalled() {
		return fmt.Errorf("migrate is not installed. Please install Migrate to continue")
	}

	missingDBs := []string{}

	for _, dbName := range dbNames {
		exists, err := checkDBsExists(dbName)
		if err != nil {
			return err
		}
		if !exists {
			missingDBs = append(missingDBs, dbName)
		}
	}

	if len(missingDBs) > 0 {
		fmt.Printf("The following databases are missing: %v\n", missingDBs)
		fmt.Print("Do you want to create them? (yes/no): ")
		var response string
		fmt.Scanln(&response)

		if strings.ToLower(response) != "yes" && strings.ToLower(response) != "y" {
			return fmt.Errorf("database initialization aborted")
		}

		for _, dbName := range missingDBs {
			err := createDB(cfg, dbName)
			if err != nil {
				return err
			}

			var dbURL string
			if dbName == "gitserver" {
				dbURL = cfg.DatabaseURL
			} else if dbName == "gitserver_test" {
				dbURL = cfg.TestDatabaseURL
			} else {
				return fmt.Errorf("unknown database name: %s", dbName)
			}

			if err := runMigrations(dbURL); err != nil {
				return fmt.Errorf("failed to run migrations on database %s: %v", dbName, err)
			}
		}
	}

	if err := validateConfigDBData(cfg); err != nil {
		return err
	}

	return nil
}
