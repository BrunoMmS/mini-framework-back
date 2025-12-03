package server

import (
	httperr "glac/custom_errors"
	"glac/router"
	"net/http"
	"io"
	"fmt"
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
	defer func() {
        if rec := recover(); rec != nil {
            s.Logger.LogError(r.Method, r.URL.Path, fmt.Errorf("%v", rec), r.RemoteAddr)
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        }
    }()

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
		Req:    r,
		Body:   bodyBytes,

		Params: make(map[string]string),
		Query:  make(map[string]string),
	}

	for k, v := range params {
		ctx.Params[k] = v
	}

	handler(ctx)
	s.Logger.LogRequest(method, path, ip)
}
