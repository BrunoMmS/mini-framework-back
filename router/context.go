package router

import (
	"encoding/json"
	"fmt"
	"net"
)

type Context struct {
	Method string
	Path   string
	Conn   net.Conn
	Body   []byte
	Params map[string]string
	Query  map[string]string
}

func (c *Context) Text(status int, body string) {
	c.writeResponse(status, "text/plain", []byte(body))
}


func (c *Context) JSON(status int, data any) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		c.Error(500, "internal server error")
		return
	}
	c.writeResponse(status, "application/json", jsonBytes)
}


func (c *Context) Error(status int, message string) {
	resp := map[string]string{"error": message}
	bytes, _ := json.Marshal(resp)
	c.writeResponse(status, "application/json", bytes)
}


func (c *Context) writeResponse(status int, contentType string, body []byte) {
	statusText := statusMessage(status)

	headers := fmt.Sprintf(
		"HTTP/1.1 %d %s\r\n"+
			"Content-Type: %s\r\n"+
			"Content-Length: %d\r\n"+
			"Connection: close\r\n"+
			"\r\n",
		status,
		statusText,
		contentType,
		len(body),
	)

	c.Conn.Write([]byte(headers))
	c.Conn.Write(body)
	c.Conn.Close()
}

func statusMessage(code int) string {
	switch code {
	case 200:
		return "OK"
	case 201:
		return "Created"
	case 400:
		return "Bad Request"
	case 404:
		return "Not Found"
	case 405:
		return "Method Not Allowed"
	case 500:
		return "Internal Server Error"
	default:
		return "OK"
	}
}
