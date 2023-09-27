package gorouter

import "net/http"

type HandleFunc func(http.ResponseWriter, *http.Request)

type Router struct {
	routes map[string]HandleFunc
}

func NewRouter() *Router {
	return &Router{
		routes: make(map[string]HandleFunc),
	}
}

func (r *Router) AddRoute(method, path string, handler HandleFunc) {
	key := method + "_" + path

	r.routes[key] = handler
}

func (r *Router) GET(path string, handler HandleFunc) {
	key := "GET_" + path
	r.routes[key] = handler
}

func (r *Router) POST(path string, handler HandleFunc) {
	key := "POST_" + path
	r.routes[key] = handler
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key := req.Method + "_" + req.URL.Path
	if handler, ok := r.routes[key]; ok {
		handler(w, req)
	} else {
		http.NotFound(w, req)
	}
}
