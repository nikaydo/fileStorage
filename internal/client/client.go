package client

import (
	"encoding/binary"
	"fmt"
	"net"
)

type Client struct {
	Port string
	Addr string
	Conn net.Conn
}

func Connect(port, addr string) (Client, error) {
	c := Client{Port: port, Addr: addr}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", c.Addr, c.Port))
	if err != nil {
		return c, err
	}
	c.Conn = conn
	return c, nil
}

func (c *Client) Write(file string) error {
	nameBytes := []byte(file)
	if len(nameBytes) > 65535 {
		return fmt.Errorf("filename too long")
	}
	err := binary.Write(c.Conn, binary.BigEndian, uint64(len(nameBytes)))
	if err != nil {
		return err
	}

	_, err = c.Conn.Write(nameBytes)
	if err != nil {
		return err
	}
	return nil
}
