package web

import (
	"fmt"
	"log"
	"strings"
)

type Node struct {
	path       string
	curPart    string
	children   []*Node
	indistinct bool
	count      int
	//pathParams map[string]string
}

func (n Node) String() string {
	if n.children == nil {
		return fmt.Sprintf("path: %s, curPart: %s, children: [], indistinct: %t, count: %d", n.path, n.curPart, n.indistinct, n.count)
	} else {
		tmp := "["
		for i := range n.children {
			tmp += (*n.children[i]).String()
			tmp += ", "
		}
		tmp += "]"
		return fmt.Sprintf("path: %s, curPart: %s, children: %s, indistinct: %t, count: %d", n.path, n.curPart, tmp, n.indistinct, n.count)
	}
}

type IRouter interface {
	Get(string, HandleFunc)
	Post(string, HandleFunc)
	Use(middleware Middleware) *IRouter
	Group(path string, handler ...HandleFunc) *RouterGroup
}

type Router struct {
	roots    map[string]*Node
	handlers map[string]HandleFunc
}

type RouterGroup struct {
	prefix      string
	children    []*RouterGroup
	middlewares []Middleware
	engine      *Engine
}

func (r *Router) insert(method string, path string, handler HandleFunc) {
	//key := method + "||" + path
	//r.handlers[key] = handler
	root, ok := r.roots[method]
	if !ok {
		root = &Node{
			path:       "",
			curPart:    "",
			children:   nil,
			indistinct: false,
			count:      0,
		}
		r.roots[method] = root
	}
	patterns := strings.Split(path, "/")
	patterns = patterns[1:]

	lastNode, longestMatch := root.search(0, patterns)
	fmt.Println(patterns)
	fmt.Printf("%+v \n", *lastNode)
	fmt.Println(longestMatch)
	if longestMatch == len(patterns) {
		key := method + "||" + path
		r.handlers[key] = handler
		return
	}
	if lastNode != nil {
		curNode := lastNode
		for i := longestMatch + 1; i < len(patterns); i++ {

			part := patterns[i]
			newNode := &Node{
				path:       "",
				curPart:    part,
				children:   nil,
				indistinct: part == "*" || strings.HasPrefix(part, ":"),
				count:      curNode.count + 1,
			}
			//curNode.children
			if curNode.children == nil {
				curNode.children = make([]*Node, 0)
			}
			curNode.children = append(curNode.children, newNode)
			curNode = newNode
		}
		curNode.path = path
		key := method + "||" + path
		r.handlers[key] = handler
		return
	}
}

//根据对应method和path查找对应的node
func (r *Router) find(method string, path string) (HandleFunc, map[string]string) {
	patterns := strings.Split(path, "/")[1:]
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}

	res, _ := root.search(0, patterns)
	fmt.Printf("FOUND ROUTER: %s", (*res).String())
	if res.path != "" {
		if res.curPart != "*" && res.count != len(patterns) {
			return nil, nil
		}
		key := method + "||" + res.path

		handler, ok := r.handlers[key]
		if !ok {
			return nil, nil
		}
		routePatterns := strings.Split(res.path, "/")[1:]
		params := map[string]string{}
		for i := range routePatterns {
			if strings.HasPrefix(routePatterns[i], ":") {
				param := string([]rune(routePatterns[i])[1:])
				params[param] = patterns[i]
			}
		}
		return handler, params
	}
	return nil, nil
}

//找到第一个与pattern匹配的node(比较index和len(pattern)来确定是否找到),否则返回最长匹配node
func (root *Node) search(index int, patterns []string) (*Node, int) {
	var longestNode = root
	maxDepth := index - 1
	if index >= len(patterns) {
		return root, len(patterns)
	}
	pattern := patterns[index]
	//插入新route 且为通配符
	if pattern == "*" {
		return root, index
	}
	if root.curPart == "*" {
		longestNode = root
	}

	for i := range root.children {
		n := root.children[i]
		if n.curPart == patterns[index] || n.indistinct {
			//param := string([]rune(n.curPart)[1:len(n.curPart)])
			next, depth := n.search(index+1, patterns)
			if next != nil && (depth > maxDepth || depth == maxDepth && longestNode.indistinct) {
				longestNode = next
				maxDepth = depth
			}
		}
	}
	//if longestNode == nil || (longestNode.curPart != "*" && maxDepth != len(pattern)){
	//	return nil, -1
	//}

	return longestNode, maxDepth

}

func (r *RouterGroup) Group(path string, handler ...HandleFunc) *RouterGroup {
	//e.router.addRoute("GET", path, handler)
	newGroup := &RouterGroup{
		prefix: r.prefix + path,
		engine: r.engine,
	}
	r.engine.children = append(r.engine.children, newGroup)
	return newGroup
}

func (r *RouterGroup) addRoute(method, path string, handler HandleFunc) *RouterGroup {
	//e.router.addRoute("GET", path, handler)
	pattern := r.prefix + path
	log.Printf("Route %4s - %s", method, pattern)
	r.engine.router.addRoute(method, pattern, handler)
	return r
}

func (r *Router) addRoute(method string, path string, handler HandleFunc) {
	r.insert(method, path, handler)
}

func (r *Router) getHandler(method, path string) (HandleFunc, map[string]string) {
	return r.find(method, path)
}
