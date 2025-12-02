package router

import (
		httperr "glac/errors"
	)


type HandlerFunc func() (any, error)

type Node struct {
	Function HandlerFunc
	Method string
}

type Router struct {
	Routes map[string]*Node
}


func (r *Router) Handle(method string, path string, fn HandlerFunc){
	newNode := &Node{
		Function: fn,
		Method: method,
	}
	
	r.Routes[path] = newNode
}

func (r *Router) Resolve(method string, path string) (any, error) {
    n, ok := r.Routes[path]
    if !ok {
        return nil, &httperr.HttpError{
            Status:  404,
            Message: "Not Found",
        }
    }

    if n.Method != method {
        return nil, &httperr.HttpError{
            Status:  405,
            Message: "Method Not Allowed",
        }
    }
    return n.Function()
}


func InitRouter() *Router {
	return &Router{Routes: make(map[string]*Node)}
}