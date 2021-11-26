package main

import (
	"log"
	"net-cat/server"
)

func main() {
	s := server.NewServer()
	go s.Serve()
	err := s.Listen()
	if err != nil {
		log.Fatal(err)
	}
}
