package main

import (
	"fmt"
	"log"
	"net-cat/server"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("[USAGE]: ./TCPChat $port")
	}
	port := os.Args[1]
	if _, err := strconv.Atoi(port); err != nil {
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	}
	s := server.NewServer(port)
	go s.Serve()
	err := s.Listen()
	if err != nil {
		log.Fatal(err)
	}
}
