package server

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	EnvLocal = "local"
	EnvDev   = "dev"
	EnvProd  = "prod"
)

type Config struct {
	Env  string `yaml:"env" env:"ENV,required" envDefault:"local"`
	Port string `yaml:"address" env:"ADDRESS" envDefault:"8085"`
	DB   struct {
		Host     string `yaml:"host" env:"DB_HOST" envDefault:"localhost"`
		Port     string `yaml:"port" env:"DB_PORT" envDefault:"5432"`
		User     string `yaml:"user" env:"DB_USER" envDefault:"postgres"`
		Password string `yaml:"password" env:"DB_PASSWORD" envDefault:"pass"`
		Name     string `yaml:"name" env:"DB_NAME" envDefault:"gitserver"`
		TestName string `yaml:"test_name" env:"DB_TEST_NAME" envDefault:"gitserver_test"`
	} `yaml:"db"`
	DatabaseURL     string `env:"DATABASE_URL"`
	TestDatabaseURL string `env:"TEST_DATABASE_URL"`
}

// TODO : change initialization configPath
var configPath string = "/home/timno/Documents/GoProjects/git-server/internal/configs/local.yaml"

func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig(configPath, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to read environment variables: %w", err)
	}

	if cfg.DatabaseURL == "" {
		cfg.DatabaseURL = generateDBURL(cfg, cfg.DB.Name)
	}
	if cfg.TestDatabaseURL == "" {
		cfg.TestDatabaseURL = generateDBURL(cfg, cfg.DB.TestName)
	}

	return cfg, nil
}

func generateDBURL(cfg *Config, dbName string) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DB.User, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, dbName)
}
