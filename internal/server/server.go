package server

import (
	"encoding/binary"
	"io"
	"net"
)

type Server struct {
	Port string
	Addr string
	Serv net.Listener
	Conn net.Conn
}

func ServerInit() (Server, error) {
	var s Server
	listener, err := net.Listen("tcp", s.Port)
	if err != nil {
		return s, err
	}
	s.Serv = listener
	return s, nil
}

func (s *Server) Run(f func(c net.Conn)) {
	for {
		conn, err := s.Serv.Accept()
		if err != nil {

		}
		go f(conn)
	}
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
