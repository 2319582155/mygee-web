package gee

import (
	"fmt"
	"log"
	"strings"
)

type router struct {
	roots    map[string]*node
	handlers map[string]HandleFunc
}

func newRoute() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandleFunc),
	}
}

func parsePattern(pattern string) []string {
	split := strings.Split(pattern, "/")
	parts := make([]string, 0)
	for _, item := range split {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *router) addRoute(method string, pattern string, handle HandleFunc) {
	parts := parsePattern(pattern)
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	key := method + "-" + pattern
	r.handlers[key] = handle
}

func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	n, ok := r.roots[method]
	res := make(map[string]string)

	if !ok {
		return nil, nil
	}

	n = n.search(searchParts, 0)

	if n != nil {
		parts := parsePattern(n.pattern)
		for i, part := range parts {
			if part[0] == ':' {
				res[part[1:]] = searchParts[i]
			}
			if part[0] == '*' && len(part) > 1 {
				res[part[1:]] = strings.Join(searchParts[i:], "/")
			}
		}
		return n, res
	}
	return nil, nil
}

func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		c.params = params
		key := c.Method + "-" + n.pattern
		handle, ok := r.handlers[key]
		if !ok {
			log.Println("没有找到路由对应的处理器，请检查路由映射是否正确，key是否正确")
		}
		handle(c)
	} else {
		_, err := fmt.Fprintf(c.Writer, "404 NOT FOUND: %s\n", c.Request.URL)
		if err != nil {
			fmt.Println(err)
		}
	}
}
