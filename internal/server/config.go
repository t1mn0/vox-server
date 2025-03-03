package server

const (
	EnvLocal = "local"
	EnvDev   = "dev"
	EnvProd  = "prod"
)

type Config struct {
	Env         string `yaml:"env" env-default:"local" env-required:"true"`
	Port        string `yaml:"address" env-default:"8085"`
	DatabaseURL string `yaml:"database_url"`
}

func NewConfig() *Config {
	return &Config{
		Port: ":8085",
		Env:  "local",
	}
}
