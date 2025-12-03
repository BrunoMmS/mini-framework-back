package router

type MiddlewareFunc func(HandlerFunc) HandlerFunc

func (m MiddlewareFunc) Handle(next HandlerFunc) HandlerFunc {
	return m(next)
}