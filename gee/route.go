package gee

import "fmt"

type Route struct {
	handles map[string]HandleFunc
}

func newRoute() *Route {
	return &Route{handles: make(map[string]HandleFunc)}
}

func (r *Route) addRoute(method string, pattern string, handle HandleFunc) {
	key := method + "-" + pattern
	r.handles[key] = handle
}

func (r *Route) handle(c *Context) {
	key := c.Request.Method + "-" + c.Request.URL.Path
	if handle, ok := r.handles[key]; ok {
		handle(c)
	} else {
		_, err := fmt.Fprintf(c.Writer, "404 NOT FOUND: %s\n", c.Request.URL)
		if err != nil {
			fmt.Println(err)
		}
	}
}
