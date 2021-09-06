package gil

import (
	"log"
	"net/http"
	"strings"
)

// HandlerFunc 定义了 gee 框架的 handler 函数
type HandlerFunc func(*Context)

type (
	RouterGroup struct {
		prefix      string
		// 支持中间件，中间件是应用在 RouterGroup 上的
		middlewares []HandlerFunc
		// 支持分组嵌套
		parent      *RouterGroup
		// 所有路由组共享一个 Engine 实例
		engine      *Engine
	}

	// Engine 继承了 RouterGroup，并实现了 http.Handler接口
	Engine struct {
		*RouterGroup
		router *router
		// 所有的路由组
		groups []*RouterGroup
	}
)

// New 是 Engine 的构造函数
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// Group 用于创建新的路由组
// 因为 Engine 继承了 RouterGroup，所以 Engine 也有 Group 方法，下面同理
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

// Use 用来想路由组添加中间件
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

// GET 用于处理 GET 请求
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

// POST 用于处理 POST 请求
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

// DELETE 用于处理 DELETE 请求
func (group *RouterGroup) DELETE(pattern string, handler HandlerFunc) {
	group.addRoute("DELETE", pattern, handler)
}

// PUT 用于处理 PUT 请求
func (group *RouterGroup) PUT(pattern string, handler HandlerFunc) {
	group.addRoute("PUT", pattern, handler)
}

// PATCH 用于处理 PATCH 请求
func (group *RouterGroup) PATCH(pattern string, handler HandlerFunc) {
	group.addRoute("PATCH", pattern, handler)
}

// HEAD 用于处理 HEAD 请求
func (group *RouterGroup) HEAD(pattern string, handler HandlerFunc) {
	group.addRoute("HEAD", pattern, handler)
}

// 使 Engine 实现 http.Handler 接口，每次请求都会触发 ServeHTTP 方法
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// 从 RouterGroup 中提取匹配的 middleware
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	// 初始化上下文，并为上下文添加该 RouterGroup 中的 middelware
	c := newContext(w, req)
	c.handlers = middlewares
	engine.router.handle(c)
}

// Run 用于启动 http server
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}