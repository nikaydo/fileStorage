package server

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
)

type Server struct {
	Port string
	Addr string
	Serv net.Listener
	Conn net.Conn
}

func ServerInit(port, addr string) (Server, error) {
	s := Server{Port: port, Addr: addr}
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", s.Addr, s.Port))
	if err != nil {
		return s, err
	}
	s.Serv = listener
	return s, nil

}

func (s *Server) Connect() error {
	conn, err := s.Serv.Accept()
	if err != nil {
		return err
	}
	s.Conn = conn
	return nil
}

func (s *Server) Write(data []byte) error {
	if err := binary.Write(s.Conn, binary.BigEndian, uint64(len(data))); err != nil {
		return err
	}
	if _, err := s.Conn.Write(data); err != nil {
		return err
	}
	return nil
}

func (s *Server) Read() ([]byte, error) {
	var data uint64
	err := binary.Read(s.Conn, binary.BigEndian, &data)
	if err != nil {
		return nil, err
	}
	nameBuf := make([]byte, data)
	_, err = io.ReadFull(s.Conn, nameBuf)
	if err != nil {
		return nil, err
	}
	return nameBuf, nil
}

func (s *Server) WriteFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	inf, err := file.Stat()
	if err != nil {
		return err
	}
	if err := binary.Write(s.Conn, binary.BigEndian, uint64(inf.Size())); err != nil {
		return err
	}
	data := make([]byte, 4096)
	for {
		_, err := file.Read(data)
		if err == io.EOF {
			break
		}
		if _, err := s.Conn.Write(data); err != nil {
			return err
		}
	}
	return nil
}
