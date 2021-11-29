package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
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
		n, err := conn.Read(buff)
		if err != nil {
			log.Println("client error during reading")
			break
		}
		fmt.Print(string(buff[:n]))
	}
}

func write(conn net.Conn, reader *bufio.Reader) {
	name := make([]byte, 2048)
	reader.Read(name)
	name = bytes.Trim(name, "\x00")
	// name = bytes.Trim(name, "\r\n")
	// username := strings.Trim(string(name), "\r\n")
	conn.Write([]byte(name))
	for {
		msg := make([]byte, 2048)
		_, err := reader.Read(msg)
		if err != nil {
			log.Println(err)
			continue
		}
		conn.Write(bytes.Trim(msg, "\x00"))
	}
}
