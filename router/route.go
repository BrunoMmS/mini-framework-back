package router

type Route struct {
	Handler HandlerFunc
	Middlewares []Middleware
	Path string
}

