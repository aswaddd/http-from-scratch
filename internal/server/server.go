package server

import (
	"fmt"
	"io"
	"net"

	"aswad.http.module/internal/request"
	"aswad.http.module/internal/response"
)

type Handler func(w *response.Writer, req *request.Request)

// type HandlerFunc func(w http.ResponseWriter, r *http.Request)

type Server struct {
	closed   bool
	handler  Handler
	listener net.Listener
}

func runConnection(s *Server, conn io.ReadWriteCloser) {
	defer conn.Close()

	responseWriter := response.NewWriter(conn)
	r, err := request.RequestFromReader(conn)
	if err != nil {
		_ = responseWriter.WriteStatusLine(response.StatusBadRequest)
		defaultHeaders := response.GetDefaultHeaders(0)
		_ = responseWriter.WriteHeaders(*defaultHeaders)
		return
	}

	if s.handler != nil {
		s.handler(responseWriter, r)
	}
}

func runServer(s *Server, listener net.Listener) {
	for {
		if s.closed {
			return
		}

		conn, err := listener.Accept()
		if err != nil {
			return
		}
		go runConnection(s, conn)
	}
}

func Serve(port uint16, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	server := &Server{closed: false, handler: handler, listener: listener}
	go runServer(server, listener)

	return server, nil
}

func (s *Server) Close() error {
	s.closed = true
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}
