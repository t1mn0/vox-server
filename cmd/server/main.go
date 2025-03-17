package main

import (
	"log"
	"vox-server/internal/server"
)

func main() {
	s, err := server.StartServer(false)
	if err != nil {
		log.Fatal(err)
	}

	if err := s.RunServer(); err != nil {
		log.Fatal(err)
	}
}
