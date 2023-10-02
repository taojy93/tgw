package tgw

import (
	"log"
	"net/http"
	"strings"
)

type HandlerFunc func(*Context)

type (
	// 路由器（路由器作为根路由组）
	engine struct {
		*routerGroup
		router *router
		groups []*routerGroup
	}

	routerGroup struct {
		prefix      string
		middlewares []HandlerFunc // 该路由组的中间件
		parent      *routerGroup  // 该路由组的上级
		engine      *engine
	}
)

// http.Handler interface 的实现
func (engine *engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, req)
	c.handlers = middlewares
	engine.router.handle(c)
}

// init router engine without middleware
func New() *engine {
	engine := &engine{router: newRouter()}
	engine.routerGroup = &routerGroup{engine: engine}
	engine.groups = []*routerGroup{engine.routerGroup}
	return engine
}

// init router engine with middleware
func Default() *engine {
	engine := New()
	engine.Use(Recovery(), Logger())
	return engine
}

func (group *routerGroup) Group(prefix string) *routerGroup {
	engine := group.engine
	newGroup := &routerGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

// 加中间件
func (group *routerGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

func (group *routerGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

func (group *routerGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

func (group *routerGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

func (engine *engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}
