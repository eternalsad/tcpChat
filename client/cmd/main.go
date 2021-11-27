package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"

	"github.com/jroimartin/gocui"
)

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

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
	reader := bufio.NewReader(conn)
	for {
		incomingMsg, err := reader.ReadString(':')
		if err == io.EOF {
			conn.Close()
			fmt.Println("Connection closed")
			os.Exit(0)
		}
		fmt.Printf("%s", incomingMsg) //------------------------------\n
	}
}

func write(conn net.Conn, reader *bufio.Reader) {
	username, err := reader.ReadString(':')
	logFatal(err)
	username = strings.Trim(username, " \r\n")
	conn.Write([]byte(username))
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		msg = fmt.Sprintf("[%s]: %s\n", username, strings.Trim(msg, " \r\n"))
		conn.Write([]byte(msg))
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("hello", maxX/2-7, maxY/2, maxX/2+7, maxY/2+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "Hello world!")
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
