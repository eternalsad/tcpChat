package server

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net-cat/models"
	"os"
	"strings"
	"sync"
	"time"
)

const MaxClientAmmount = 10
const timeFormat = "2006-01-02 15:04:05"

type Server struct {
	Clients         []*models.Client
	msgs            chan *models.Message
	clientsChan     chan *models.Client
	inactiveClients chan *models.Client
	Mtx             sync.Mutex
	UserCount       int
	writer          *log.Logger
	logs            *os.File
	logsFilename    string
	port            string
}

func NewServer(port string) *Server {
	file, err := os.OpenFile("logs.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)

	if err != nil {
		log.Fatal("can't create message log file")
	}

	return &Server{
		Clients:         make([]*models.Client, 0),
		msgs:            make(chan *models.Message),
		clientsChan:     make(chan *models.Client),
		inactiveClients: make(chan *models.Client),
		UserCount:       0,
		writer:          log.New(file, "", 0),
		logsFilename:    "logs.txt",
		port:            port,
		// logs:            file,
		// reader:          bufio.NewScanner(file),
	}
}

func (s *Server) Listen() error {
	ln, err := net.Listen("tcp", ":"+s.port)
	if err != nil {
		return err
	}
	for {
		conn, err := ln.Accept()
		// conn.Write([]byte("asd"))
		if err != nil {
			log.Print(err)
			continue
		}
		client, err := s.Accept(conn)
		if err != nil {
			log.Print(err)
			continue
		}
		s.clientsChan <- client
	}
	return nil
}

func (s *Server) handleConnection(client *models.Client) {
	// client.Connection.Write([]byte("[ENTER YOUR NAME]:\n"))
	s.Welcome(client)
	rd := bufio.NewReader(client.Connection)
	name, err := rd.ReadString('\n')
	if err != nil {
		fmt.Println("error while reading name")
	}
	name = strings.Trim(name, "\r\n")
	// log.Println(name)
	client.Name = name
	// s.msgs <- &models.Message{Text: fmt.Sprintf("%v joined the chat\n", name), Source: client}
	s.loadPreviousMessages(client)
	s.Mtx.Lock()
	dt := time.Now()
	timeStamp := dt.Format(timeFormat)
	go s.BroadCast(&models.Message{Text: fmt.Sprintf("\n[%v]%v joined the chat\n", timeStamp, name), Source: client})
	s.Clients = append(s.Clients, client)
	s.UserCount++
	s.Mtx.Unlock()
	// s.BroadCast(&models.Message{Text: fmt.Sprintf("%v entered chat", name), })
	for {
		dt := time.Now()
		timeStamp := dt.Format(timeFormat)
		prefix := fmt.Sprintf("[%v][%v]:", timeStamp, name)
		client.Connection.Write([]byte(prefix))
		str, err := rd.ReadString('\n')
		fmt.Printf("message length: %v\n", len(str))
		if err != nil {
			fmt.Println("error while reading from connection")
			fmt.Println(err.Error())
			break
		}
		if str != "\n" {
			msg := &models.Message{
				Text:   str,
				Source: client,
			}
			s.msgs <- msg
		}
	}
	s.inactiveClients <- client
}

func (s *Server) Accept(conn net.Conn) (*models.Client, error) {
	s.Mtx.Lock()
	if s.UserCount < MaxClientAmmount {
		s.Mtx.Unlock()
		return &models.Client{Connection: conn, ID: s.UserCount}, nil
	}
	return nil, fmt.Errorf("Maximum ammount of users reached")
}

func (s *Server) SendMessage(msg *models.Message) {
	dt := time.Now()
	timeStamp := dt.Format(timeFormat)
	message := fmt.Sprintf("\n[%v][%v]:%v", timeStamp, msg.Source.Name, msg.Text)
	msg.Text = message
	s.BroadCast(msg)
	// s.writeToFile(message)
}

// send string message to other clients
func (s *Server) BroadCast(msg *models.Message) {
	dt := time.Now()
	timeStamp := dt.Format(timeFormat)
	s.Mtx.Lock()
	for _, c := range s.Clients {
		if c.Connection != msg.Source.Connection {
			postfix := fmt.Sprintf("[%v][%v]:", timeStamp, c.Name)
			message := msg.Text + postfix
			c.Connection.Write([]byte(message))
			log.Print(message)
		}
	}
	s.Mtx.Unlock()
	s.writeToFile(msg.Text[1:])
}

func (s *Server) writeToFile(line string) {
	// log.Println(len(line))
	s.Mtx.Lock()
	s.writer.Output(2, line)
	s.Mtx.Unlock()
	// s.writer.Printf("%v", line)
	// fmt.Fprintf(s.logs, "%v", line[1:])
}

func (s *Server) removeClient(client *models.Client) {
	client.Connection.Close()
	dt := time.Now()
	timeStamp := dt.Format(timeFormat)
	message := fmt.Sprintf("\n[%v]:%v%v", timeStamp, client.Name, " has left the chat...\n")
	// s.msgs <- &models.Message{Text: message, Source: client}
	s.Mtx.Lock()
	go s.BroadCast(&models.Message{Text: message, Source: client})
	s.Clients = append(s.Clients[:client.ID], s.Clients[client.ID:]...)
	s.UserCount--
	s.Mtx.Unlock()
}

func (s *Server) Serve() {
	for {
		select {
		case client := <-s.clientsChan:
			go s.handleConnection(client)
		case msg := <-s.msgs:
			go s.SendMessage(msg)
			// s.writeToFile()
		case inactive := <-s.inactiveClients:
			go s.removeClient(inactive)
		}
	}
}

// Welcome prints Welcome prompt to new user
func (s *Server) Welcome(client *models.Client) {
	msg := "Welcome to TCP-Chat!\n         _nnnn_\n        dGGGGMMb\n       @p~qp~~qMb\n       M|@||@) M|\n       @,----.JM|\n      JS^\\__/  qKL\n     dZP        qKRb\n    dZP          qKKb\n   fZP            SMMb\n   HZM            MMMM\n   FqM            MMMM\n __| \".        |\\dS\"qML\n |    `.       | `' \\Zq\n_)      \\.___.,|     .'\n\\____   )MMMMMP|   .'\n     `-'       `--'\n"
	msg = msg + "[ENTER YOUR NAME]: "
	client.Connection.Write([]byte(msg))
}

// loadPreviousMessage sends old messages to client
func (s *Server) loadPreviousMessages(client *models.Client) {
	// TO DO POSSIBLE MEMORY LEAK
	file, err := os.OpenFile(s.logsFilename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)

	if err != nil {
		log.Fatal("can't create message log file")
	}
	s.Mtx.Lock()
	text, err := ioutil.ReadAll(file)
	s.Mtx.Unlock()
	if err != nil {
		fmt.Print("error during loading messages")
		fmt.Println(err)
		return
	}
	client.Connection.Write(text)
}
