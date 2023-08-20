package core

import (
	"errors"
	"fmt"
	"io"
	"net"
)

type ServerOpts struct {
	Port int
}

type Server struct {
	opts             ServerOpts
	requestProcessor RequestProcessor
}

func NewServer(opts ServerOpts, requestProcessor RequestProcessor) *Server {
	return &Server{opts: opts, requestProcessor: requestProcessor}
}

func (s *Server) Start() {
	fmt.Println("Starting server on port", s.opts.Port)
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", s.opts.Port))

	if err != nil {
		fmt.Printf("Error is %s\n", err)
		return
	}

	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println("Connection closed")
				continue
			}
			fmt.Printf("Error is %s\n", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Incoming connection")

	for {
		buff := make([]byte, 4096)
		n, err := conn.Read(buff)

		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println("Connection closed")
				return
			}
			fmt.Printf("Error is %s\n", err)
			return
		}

		data_recv := string(buff[:n])
		data_recv = data_recv[:len(data_recv)-1]

		fmt.Printf("Received %d bytes %s\n", n, data_recv)

		req, err := parseProtocol(data_recv)

		if err != nil {
			fmt.Printf("Error parsing protocol %s\n", err)
			s.writeError(err, conn)
			continue
		}

		resp, err := s.requestProcessor.Process(req)

		if err != nil {
			fmt.Printf("Error parsing protocol %s\n", err)
			s.writeError(err, conn)
			continue
		}

		s.writeSuccess(resp.Value, conn)
	}
}

func (s *Server) writeError(err error, conn net.Conn) {
	conn.Write([]byte(fmt.Sprintf("ERR: %s\n", err)))
}

func (s *Server) writeSuccess(value []byte, conn net.Conn) {
	if len(value) == 0 {
		conn.Write([]byte("OK\n"))
		return
	}
	conn.Write([]byte(fmt.Sprintf("%s\n", value)))
}
