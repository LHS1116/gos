package web

import (
	"fmt"
	"net/http"
	"strconv"
)

//type HandleFunc func(*http.Request, http.ResponseWriter)
type HandleFunc func(*GosContext)

type Middleware func() HandleFunc

type Engine struct {
	//request *http.Request
	//response * http.Response
	//writer *http.ResponseWriter
	middlewares []Middleware
	router      *Router
}

func Default() *Engine {
	return &Engine{router: &Router{make(map[string]HandleFunc)}, middlewares: []Middleware{}}
}

func (e *Engine) Get(path string, handler HandleFunc) {
	e.router.addRoute("GET", path, handler)
}

func (e *Engine) Post(path string, handler HandleFunc) {
	e.router.addRoute("POST", path, handler)
}

func (e *Engine) serve(ctx *GosContext) {
	r := ctx.Request
	key := r.Method + ":" + r.URL.Path

	if handler, ok := e.router.getHandler(key); ok {
		handler(ctx)
	} else {
		_, err := fmt.Fprintf(ctx.Writer, "404 NOT FOUND: %s\n", r.URL)
		if err != nil {
			return
		}
	}
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	context := newContext(r, w, e.middlewares)
	context.Next()

	//for i, _ := range e.middlewares {
	//	m := e.middlewares[i]
	//	m()(context)
	//
	//}

	e.serve(context)
}

func (e *Engine) Use(middleware Middleware) *Engine {
	e.middlewares = append(e.middlewares, middleware)
	return e
}

func (e *Engine) Run(port int) error {
	return http.ListenAndServe(":"+strconv.Itoa(port), e)
}
