package router

import (
	"net/http"

	"gitlab.momoso.com/cm/kit/third_party/lg"
)

type middleware func(http.Handler) http.Handler

type Router struct {
	middlewareChain []middleware
	mux             map[string]http.Handler
}

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) Use(m middleware) {
	r.middlewareChain = append(r.middlewareChain, m)
}

func (r *Router) Add(route string, h http.Handler) {
	var mergedHandler = h
	for i := len(r.middlewareChain) - 1; i >= 0; i-- {
		mergedHandler = r.middlewareChain[i](mergedHandler)
	}
	if r.mux == nil {
		r.mux = make(map[string]http.Handler)
	}
	r.mux[route] = mergedHandler
}

func (r *Router) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	lg.Infof("req path:%s", path)
	if handler, ok := r.mux[path]; ok {
		lg.Infof("path matched")
		handler.ServeHTTP(rw, req)
		lg.Infof("router serve done")
	}
}
