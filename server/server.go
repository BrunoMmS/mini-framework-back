package server

import (
	"bufio"
	"fmt"
	httperr "glac/errors"
	"glac/router"
	"net"
	"strings"
)


type Server struct{
	protocol string
	port string
	Router *router.Router
}

func (s *Server) Listen() {
	ln, err := net.Listen(s.protocol, s.port)
	logger := InitLogger(INFO)
	if err != nil {
		logger.LogError("LISTEN", s.port, err, "system")
		return
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			logger.LogError("ACCEPT", "-", err, "system")
			continue
		}

		go s.handleConnection(conn, logger)
	}
}

func (s *Server) handleConnection(conn net.Conn, logger *Logger) {
	defer conn.Close()

	ip := conn.RemoteAddr().String()
	reader := bufio.NewReader(conn)

	line, err := reader.ReadString('\n')
	if err != nil {
		logger.LogError("READ", "-", err, ip)
		return
	}

	parts := strings.Split(strings.TrimSpace(line), " ")
	if len(parts) < 2 {
		logger.LogError("PARSE", "-", fmt.Errorf("invalid request line"), ip)
		conn.Write([]byte("HTTP/1.1 400 Bad Request\r\nConnection: close\r\n\r\n"))
		return
	}

	method := parts[0]
	path := parts[1]

	handler, params, err := s.Router.Resolve(method, path)
	if err != nil {
		httpErr := err.(*httperr.HttpError)
		raw := fmt.Sprintf(
			"HTTP/1.1 %d %s\r\nConnection: close\r\n\r\n",
			httpErr.Status, httpErr.Message,
		)
		conn.Write([]byte(raw))
		logger.LogError(method, path, err, ip)
		return
	}

	ctx := &router.Context{
		Method: method,
		Path:   path,
		Conn:   conn,
		Params: params,
		Body:   nil,
	}

	handler(ctx)

	logger.LogRequest(method, path, ip)
}


func NewServer(protocol string, port string, r *router.Router) *Server{
	return &Server{
		protocol: protocol,
		port: port,
		Router: r,
	}
}