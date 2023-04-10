package gee

import (
	"net/http"
)

type HandleFunc func(c *Context)

type Engine struct {
	route *Route
}

func New() *Engine {
	return &Engine{route: newRoute()}
}

func (e *Engine) Get(pattern string, handle HandleFunc) {
	e.addRoute("GET", pattern, handle)
}

func (e *Engine) Post(pattern string, handle HandleFunc) {
	e.addRoute("POST", pattern, handle)
}

func (e *Engine) addRoute(method string, pattern string, handle HandleFunc) {
	e.route.addRoute(method, pattern, handle)
}

func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := newContext(w, r)
	e.route.handle(c)
}
