package gil

import (
	"net/http"
	"strings"
)

// 封装了前缀树
type router struct {
	// 存储每种请求方式（GET, POST 等）的前缀树根节点
	roots    map[string]*node
	// 存储每个请求（handlers['GET-/p/:lang/doc']）对应的 HandlerFunc
	// handlers 由 router 保存，而不是由 Engine 保存
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// 将待匹配路由按 / 分割为不同部分，遇到 * 就结束
func parsePattern(pattern string) []string {
	splits := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range splits {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *router) addRoute(method string, path string, handler HandlerFunc) {
	parts := parsePattern(path)

	key := method + "-" + path
	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(path, parts, 0)
	r.handlers[key] = handler
}

func (r *router) getRoute(method string, pattern string) (*node, map[string]string) {
	searchParts := parsePattern(pattern)
	params := make(map[string]string)
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}

	node := root.search(searchParts, 0)
	if node != nil {
		parts := parsePattern(node.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return node, params
	}

	return nil, nil
}

func (r *router) getRoutes(method string) []*node {
	root, ok := r.roots[method]
	if !ok {
		return nil
	}
	nodes := make([]*node, 0)
	root.travel(&nodes)
	return nodes
}

// 设置 context 的 handlers 和 params，并调用 handler 
func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		key := c.Method + "-" + n.pattern
		c.handlers = append(c.handlers, r.handlers[key])
		c.Params = params
	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}
	// 开启 handler chain
	c.Next()
}
