package gee

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	Writer     http.ResponseWriter
	Request    *http.Request
	Method     string
	Path       string
	params     map[string]string
	StatusCode int
}

func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer:  w,
		Request: r,
		Method:  r.Method,
		Path:    r.URL.Path,
	}
}

func (c *Context) Query(key string) string {
	return c.Request.URL.Query().Get(key)
}

func (c *Context) PostForm(key string) string {
	return c.Request.FormValue(key)
}

func (c *Context) status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Context) setHeader(key string, val string) {
	c.Writer.Header().Set(key, val)
}

func (c *Context) JSON(code int, data interface{}) {
	c.setHeader("Content-type", "application/json")
	c.status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(data); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
	}
}

func (c *Context) String(code int, format string, values ...interface{}) {
	c.setHeader("Content-type", "text/plain")
	c.status(code)
	_, err := c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
	if err != nil {
		log.Println(err)
	}
}

func (c *Context) Data(code int, data []byte) {
	c.status(code)
	log.Println(c.Writer.Write(data))
}

func (c *Context) HTML(code int, html string) {
	c.setHeader("Content-type", "text/html")
	c.status(code)
	_, err := c.Writer.Write([]byte(html))
	if err != nil {
		fmt.Println(err)
	}
}

func (c *Context) Param(key string) string {
	value := c.params[key]
	return value
}
