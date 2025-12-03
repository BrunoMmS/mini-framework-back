package router

type Middleware func(HandlerFunc) HandlerFunc

func Chain(mws []Middleware, final HandlerFunc) HandlerFunc {
    if len(mws) == 0 {
        return final
    }

    wrapped := final
    for i := len(mws) - 1; i >= 0; i-- {
        wrapped = mws[i](wrapped)
    }
    return wrapped
}
