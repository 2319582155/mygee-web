package gee

import (
	"log"
	"net/http"
)

type HandleFunc func(c *Context)

// Engine 整个框架的所有资源由Engine统一协调
type Engine struct {
	*RouterGroup
	route  *router
	groups []*RouterGroup
}

// RouterGroup 分组结构，进行分组控制
type RouterGroup struct {
	prefix      string       // prefix : 分组前缀
	parent      *RouterGroup // parent : 父亲
	middlewares []HandleFunc // middlewares : 给中间件提供无线扩展能力的各种中间件
	engine      *Engine      // engine : 所有的Group共享一个资源引擎Engine
}

func New() *Engine {
	engine := &Engine{route: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) Get(pattern string, handle HandleFunc) {
	group.addRoute("GET", pattern, handle)
}

func (group *RouterGroup) Post(pattern string, handle HandleFunc) {
	group.addRoute("POST", pattern, handle)
}

func (group *RouterGroup) addRoute(method string, comp string, handle HandleFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.route.addRoute(method, pattern, handle)
}

func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := newContext(w, r)
	e.route.handle(c)
}
