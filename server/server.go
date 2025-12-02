package server

import (
	"bufio"
	"fmt"
	"glac/router"
	"net"
	"strings"
)


type Server struct{
	protocol string
	port string
	Router *router.Router
}

func (s *Server) Listen(){
	ln, err := net.Listen(s.protocol, s.port)
	
	if err != nil{
		fmt.Println(err)
		return
	}

	for {
		conn, err := ln.Accept()
		if err != nil{
			fmt.Println(err)
			continue
		}

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
    defer conn.Close()

    reader := bufio.NewReader(conn)

    line, err := reader.ReadString('\n')
    if err != nil {
        fmt.Println("Error leyendo:", err)
        return
    }

    parts := strings.Split(line, " ")
    if len(parts) < 2 {
        conn.Write([]byte("HTTP/1.1 400 Bad Request\r\n\r\n"))
        return
    }

    method := parts[0]
    path := parts[1]

    res, err := s.Router.Resolve(method, path)
    if err != nil {
        conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n" + err.Error()))
        return
    }

    conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\n"))
    conn.Write([]byte(fmt.Sprintf("%v", res)))
}
func NewServer(protocol string, port string, r *router.Router) *Server{
	return &Server{
		protocol: protocol,
		port: port,
		Router: r,
	}
}