package router

import (
	httperr "glac/custom_errors"
	"strings"
)

type Router struct {
    static  map[string]map[string]*Route
    dynamic map[string][]*Route //migrate to trie
}


func (r *Router) handle(method, path string, fn HandlerFunc, mws ...Middleware) {
    rt := &Route{
        Handler:     fn,
        Middlewares: mws,
        Path:        path,
    }

    if strings.Contains(path, ":") {
        r.dynamic[method] = append(r.dynamic[method], rt)
    } else {
        if r.static[method] == nil {
            r.static[method] = make(map[string]*Route)
        }
        r.static[method][path] = rt
    }
}

func (r *Router) Get(path string, fn HandlerFunc, middlewares ...Middleware) {
    r.handle(GET, path, fn, middlewares...) 
}

func (r *Router) Post(path string, fn HandlerFunc, middlewares ...Middleware) {
    r.handle(POST, path, fn, middlewares...)
}

func (r *Router) Put(path string, fn HandlerFunc, middlewares ...Middleware) {
    r.handle(PUT, path, fn, middlewares...) 
}

func (r *Router) Delete(path string, fn HandlerFunc, middlewares ...Middleware) {
    r.handle(DELETE, path, fn, middlewares...)
}

func (r *Router) Resolve(method, path string) (HandlerFunc, map[string]string, error) {
	//static routes
    if routesByPath, ok := r.static[method]; ok {
        if rt := routesByPath[path]; rt != nil {
            final := Chain(rt.Middlewares, rt.Handler)
            return final, nil, nil
        }
    }

	//dynamic routes like /users/:id
    if dynRoutes, ok := r.dynamic[method]; ok {
        for _, rt := range dynRoutes {
            params, ok := matchPattern(rt.Path, path)
            if ok {
                final := Chain(rt.Middlewares, rt.Handler)
                return final, params, nil
            }
        }
    }

    // not found
    return nil, nil, &httperr.HttpError{
        Status:  404,
        Message: "Not Found",
    }
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
    return &Router{
        static:  make(map[string]map[string]*Route),
        dynamic: make(map[string][]*Route),
    }
}
