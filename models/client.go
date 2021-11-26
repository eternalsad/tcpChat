package models

import "net"

type Client struct {
	Connection net.Conn
	Name       string
	ID         int
}
