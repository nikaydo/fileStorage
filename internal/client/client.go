package client

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
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

func (c *Client) Read() ([]byte, error) {
	var data uint64
	err := binary.Read(c.Conn, binary.BigEndian, &data)
	if err != nil {
		return nil, err
	}
	nameBuf := make([]byte, data)
	_, err = io.ReadFull(c.Conn, nameBuf)
	if err != nil {
		return nil, err
	}
	return nameBuf, nil
}

func (c *Client) Write(data []byte) error {
	if err := binary.Write(c.Conn, binary.BigEndian, uint64(len(data))); err != nil {
		return err
	}
	if _, err := c.Conn.Write(data); err != nil {
		return err
	}
	return nil
}

func (c *Client) WriteFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	inf, err := file.Stat()
	if err != nil {
		return err
	}
	if err := c.Write([]byte(filepath.Base(file.Name()))); err != nil {
		return err
	}
	if err := binary.Write(c.Conn, binary.BigEndian, uint64(inf.Size())); err != nil {
		return err
	}
	data := make([]byte, 4096)
	for {
		_, err := file.Read(data)
		if err == io.EOF {
			break
		}
		if _, err := c.Conn.Write(data); err != nil {
			return err
		}
	}
	return nil
}
