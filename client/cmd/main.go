package main

import (
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/jroimartin/gocui"
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
	// reader := bufio.NewReader(os.Stdin)

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)

	err = initKeybindings(g, connection)
	if err != nil {
		logFatal(err)
	}

	// go read(connection, g)
	go func() {

		for {
			buff := make([]byte, 2048)
			_, err := connection.Read(buff)
			if err != nil {
				log.Println("client error during reading")
				break
			}
			// fmt.Print(string(buff[:n]))
			g.Update(func(g *gocui.Gui) error {
				v, err := g.View("output")
				if err != nil {
					return err
				}
				// v.Clear()
				// if len(buff)-1 > n && buff[n-1] != '\n' {
				// 	fmt.Fprint(v, string(buff[:n-1]))
				// }
				fmt.Fprint(v, string(buff))
				// fmt.Fprint(v, "\n")
				return nil
			})
		}
	}()
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

	// write(connection, reader)
}

func read(conn net.Conn, g *gocui.Gui) {
	// fmt.Fprint(conn, "\n")
	for {
		buff := make([]byte, 2048)
		n, err := conn.Read(buff)
		if err != nil {
			log.Println("client error during reading")
			break
		}
		// fmt.Print(string(buff[:n]))
		g.Update(func(g *gocui.Gui) error {
			v, err := g.View("output")
			if err != nil {
				return err
			}
			// v.Clear()
			// if len(buff)-1 > n && buff[n-1] != '\n' {
			// 	fmt.Fprint(v, string(buff[:n-1]))
			// }
			fmt.Fprint(v, string(buff[:n]))
			// fmt.Fprint(v, "\n")
			return nil
		})
	}
}

// func write(conn net.Conn, reader *bufio.Reader) {
// 	name := make([]byte, 2048)
// 	reader.Read(name)
// 	name = bytes.Trim(name, "\x00")
// 	// name = bytes.Trim(name, "\r\n")
// 	// username := strings.Trim(string(name), "\r\n")
// 	conn.Write([]byte(name))
// 	for {
// 		msg := make([]byte, 2048)
// 		_, err := reader.Read(msg)
// 		if err != nil {
// 			log.Println(err)
// 			continue
// 		}
// 		conn.Write(bytes.Trim(msg, "\x00"))
// 	}
// }

func write(conn net.Conn, text []byte) {
	//name := make([]byte, 2048)
	//name = []byte("kek")
	// name = bytes.Trim(name, "\x00")
	// name = bytes.Trim(name, "\r\n")
	// username := strings.Trim(string(name), "\r\n")
	//conn.Write(name)
	_, err := conn.Write(text)
	if err != nil {
		return
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	msgs, err := g.SetView("output", 0, 0, maxX-10, 25)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		if _, err := g.SetCurrentView("output"); err != nil {
			return err
		}
		msgs.Title = "Chat room"
		msgs.Autoscroll = true
	}

	input, err := g.SetView("input", 20, maxY-10, maxX-24, maxY-1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		if _, err := g.SetCurrentView("input"); err != nil {
			return err
		}
		input.Title = "Enter Text"
		input.Editable = true
		input.Wrap = true
	}
	return nil
}

func initKeybindings(g *gocui.Gui, conn net.Conn) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return gocui.ErrQuit
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("input", gocui.KeyEnter, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			text := make([]byte, 1024)
			n, err := v.Read(text)
			if err != nil {
				return err
			}
			// ctext := strings.TrimRight(string(text), " ")
			g.Update(func(g *gocui.Gui) error {
				fmt.Fprint(conn, string(text[:n]))
				outputView, err := g.View("output")
				if err != nil {
					return err
				}

				fmt.Fprint(outputView, string(text[:n]))
				fmt.Fprint(outputView, "\n")
				// fmt.Fprint(outputView, string("message"))
				// fmt.Fprint(conn, string(text))
				// write(conn, text[:n])
				return nil
			})
			v.Clear()
			if err := v.SetCursor(0, 0); err != nil {
				return err
			}
			return nil
		}); err != nil {
		return err
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
