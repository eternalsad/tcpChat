package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net-cat/models"
	"strings"
	"sync"
	"time"
)

const MaxClientAmmount = 10

type Server struct {
	Clients         []*models.Client
	msgs            chan *models.Message
	clientsChan     chan *models.Client
	inactiveClients chan *models.Client
	Mtx             sync.Mutex
	UserCount       int
}

func NewServer() *Server {
	return &Server{
		Clients:         make([]*models.Client, 0),
		msgs:            make(chan *models.Message),
		clientsChan:     make(chan *models.Client),
		inactiveClients: make(chan *models.Client),
		UserCount:       0,
	}
}

func (s *Server) Listen() error {
	ln, err := net.Listen("tcp", ":4000")
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
	s.Mtx.Lock()
	s.Clients = append(s.Clients, client)
	s.UserCount++
	s.Mtx.Unlock()
	rd := bufio.NewReader(client.Connection)
	client.Connection.Write([]byte("Enter your name\n"))
	name, err := rd.ReadString('\n')
	name = strings.Trim(name, "\r\n")
	client.Name = name
	if err != nil {
		fmt.Println("error while reading name")
	}
	for {
		dt := time.Now()
		timeStamp := dt.Format("2006-01-02 15:04:05")
		prefix := fmt.Sprintf("[%v][%v]:", timeStamp, name)
		client.Connection.Write([]byte(prefix))
		str, err := rd.ReadString('\n')
		if err != nil {
			fmt.Println("error while reading from connection")
			fmt.Println(err.Error())
			break
		}
		if str != "\n" {
			// msgs <- fmt.Sprintf("[%v] [%v]: %v", timeStamp, name, m)

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
	if s.UserCount < MaxClientAmmount {
		return &models.Client{Connection: conn, ID: s.UserCount}, nil
	}
	return nil, fmt.Errorf("Maximum ammount of users reached")
}

func (s *Server) BroadCast(msg *models.Message) {
	dt := time.Now()
	timeStamp := dt.Format("2006-01-02 15:04:05")
	message := fmt.Sprintf("\n[%v][%v]:%v", timeStamp, msg.Source.Name, msg.Text)

	s.Mtx.Lock()
	for _, c := range s.Clients {
		if c.Connection != msg.Source.Connection {
			postfix := fmt.Sprintf("[%v][%v]:", timeStamp, c.Name)
			c.Connection.Write([]byte(message + postfix))
		}
	}
	s.Mtx.Unlock()
}

func (s *Server) removeClient(client *models.Client) {
	client.Connection.Close()
	dt := time.Now()
	timeStamp := dt.Format("2006-01-02 15:04:05")
	message := fmt.Sprintf("\n[%v]:%v%v", timeStamp, client.Name, " has left the chat...\n")
	s.Mtx.Lock()
	s.Clients = append(s.Clients[:client.ID], s.Clients[client.ID:]...)
	s.UserCount--
	for _, c := range s.Clients {
		postfix := fmt.Sprintf("[%v][%v]:", timeStamp, c.Name)
		c.Connection.Write([]byte(message + postfix))
	}
	s.Mtx.Unlock()
}

func (s *Server) Serve() {
	for {
		select {
		case client := <-s.clientsChan:
			go s.handleConnection(client)
		case msg := <-s.msgs:
			go s.BroadCast(msg)
		case inactive := <-s.inactiveClients:
			go s.removeClient(inactive)
		}
	}
}
