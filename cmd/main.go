package main

import (
	"fmt"
	"log"
	"net-cat/server"
	"os"
	"strconv"
)

func main() {
	port := "8989"
	if len(os.Args) == 2 {
		port = os.Args[1]
	}
	if len(os.Args) > 2 {
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	}
	if _, err := strconv.Atoi(port); err != nil {
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	}
	fmt.Printf("Listening at port :%v\n", port)
	s := server.NewServer(port)
	go s.Serve()
	err := s.Listen()
	if err != nil {
		log.Fatal(err)
	}
}
