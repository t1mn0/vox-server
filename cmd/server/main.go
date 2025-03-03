package main

import (
	"flag"
	"git-server/internal/server"
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", "internal/configs/local.yaml", "path to config file")
}

func main() {
	flag.Parse()
	config := server.NewConfig()
	if err := cleanenv.ReadConfig(configPath, config); err != nil {
		log.Fatal(err)
	}

	s, err := server.NewServerWithDB(config)
	if err != nil {
		log.Fatal(err)
	}

	if err := s.RunServer(); err != nil {
		log.Fatal(err)
	}
}
