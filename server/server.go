package server

import (
	httperr "glac/errors"
	"glac/router"
	"net/http"
	"io"
)

type Server struct {
	Router *router.Router
	Logger *Logger
}

func NewServer(r *router.Router) *Server {
	return &Server{
		Router: r,
		Logger: InitLogger(INFO),
	}
}

func (s *Server) Listen(port string) error {
	s.Logger.LogRequest("INFO", port, "Starting HTTP server...")
	return http.ListenAndServe(port, s)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ip := r.RemoteAddr
	method := r.Method
	path := r.URL.Path
	bodyBytes, _ := io.ReadAll(r.Body)
	handler, params, err := s.Router.Resolve(method, path)
	if err != nil {
		httpErr := err.(*httperr.HttpError)

		http.Error(w, httpErr.Message, httpErr.Status)
		s.Logger.LogError(method, path, err, ip)
		return
	}

	ctx := &router.Context{
		Method: method,
		Path:   path,
		Writer: w,           
		Params: params,
		Body:   bodyBytes,      
		Req:    r,           
	}

	handler(ctx)
	s.Logger.LogRequest(method, path, ip)
}
