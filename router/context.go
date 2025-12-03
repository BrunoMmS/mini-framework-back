package router

import (
	"encoding/json"
	"net/http"
	"glac/custom_errors"
)

type Context struct {
	Writer http.ResponseWriter
	Req    *http.Request

	Method string
	Path   string
	Body   []byte
	Params map[string]string
	Query  map[string]string
}

func (c *Context) Text(status int, body string) {
	c.Writer.Header().Set("Content-Type", "text/plain")
	c.Writer.WriteHeader(status)
	c.Writer.Write([]byte(body))
}

func (c *Context) JSON(status int, data any) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(status)

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		http.Error(c.Writer, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	c.Writer.Write(jsonBytes)
}

func (c *Context) Error(err *custom_errors.HttpError) {
    c.JSON(err.Status, map[string]string{"error": err.Message})
}
