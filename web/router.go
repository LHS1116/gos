package web

type Router struct {
	handlers map[string]HandleFunc
}

func (r *Router) addRoute(method string, path string, handler HandleFunc) {
	key := method + ":" + path
	r.handlers[key] = handler
}

func (r *Router) getHandler(key string) (HandleFunc, bool) {
	f, ok := r.handlers[key]
	return f, ok
}
