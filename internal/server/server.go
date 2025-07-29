package server

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
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
