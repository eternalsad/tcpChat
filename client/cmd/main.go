package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
)

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var mtx sync.Mutex

func main() {
	connection, err := net.Dial("tcp", "localhost:4000")
	logFatal(err)

	defer connection.Close()
	reader := bufio.NewReader(os.Stdin)

	go read(connection)
	write(connection, reader)
}

func read(conn net.Conn) {
	for {
		buff := make([]byte, 2048)
		conn.Read(buff)
		fmt.Print(string(buff))
	}
}

func write(conn net.Conn, reader *bufio.Reader) {
	name := make([]byte, 2048)
	reader.Read(name)
	username := strings.Trim(string(name), "\n")
	conn.Write([]byte(username))
	for {
		msg := make([]byte, 2048)
		reader.Read(msg)
		conn.Write([]byte(msg))
	}
}
