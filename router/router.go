package router

import (
	httperr "glac/errors"
	"strings"
)


type HandlerFunc func(c *Context)

type Node struct {
    Method string
	Handler HandlerFunc
}

type Router struct {
	routes map[string]*Node
}


func (r *Router) Handle(method string, path string, fn HandlerFunc){
	newNode := &Node{
		Handler: fn,
		Method: method,
	}
	
	r.routes[path] = newNode
}

func (r *Router) Resolve(method, path string) (HandlerFunc, map[string]string, error) {
    for pattern, node := range r.routes {
        if node.Method != method {
            continue
        }

        params, ok := matchPattern(pattern, path)
        if ok {
            return node.Handler, params, nil
        }
    }

    return nil, nil, &httperr.HttpError{Status: 404, Message: "Not Found"}
}

func matchPattern(pattern, path string) (map[string]string, bool) {
    params := make(map[string]string)

    p1 := splitPath(pattern)
    p2 := splitPath(path)

    if len(p1) != len(p2) {
        return nil, false
    }

    for i := 0; i < len(p1); i++ {
        seg := p1[i]
        val := p2[i]

        if strings.HasPrefix(seg, ":") {
            params[seg[1:]] = val
            continue
        }

        if seg != val {
            return nil, false
        }
    }

    return params, true
}

func splitPath(p string) []string {
    parts := strings.Split(p, "/")
    out := make([]string, 0, len(parts))
    for _, v := range parts {
        if v != "" {
            out = append(out, v)
        }
    }
    return out
}

func InitRouter() *Router {
	return &Router{routes: make(map[string]*Node)}
}