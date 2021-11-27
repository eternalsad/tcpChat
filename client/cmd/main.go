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
	// color.Cyan.Print("Enter you name: ")
	reader := bufio.NewReader(os.Stdin)
	// connection.Write([]byte(username))

	go read(connection)
	write(connection, reader)
}

func read(conn net.Conn) {
	// mtx.Lock()
	rd := bufio.NewReader(conn)
	// mtx.Unlock()
	for {
		// mtx.Lock()
		msg, err := rd.ReadString(':')
		if err != nil {
			conn.Close()
			return // maybe os.Exit()?
		}
		fmt.Print(msg)
		// mtx.Unlock()
	}
}

func write(conn net.Conn, reader *bufio.Reader) {
	scanner := bufio.NewScanner(reader)

	username := scanner.Text() // ne ponimau pochemu on reagiruet na \r
	// a ne na \n maybe potomu chto tam sequence \r\n i on jdal \n
	// logFatal(err)
	username = strings.Trim(username, "\r\n ")
	conn.Write([]byte(username))
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		conn.Write([]byte(msg))
	}
}
