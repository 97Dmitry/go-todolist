package httpserver

import "net/http"

type APIVersion string

const APIVersionV1 APIVersion = "v1"

type APIVersionRouter struct {
	mux        *http.ServeMux
	apiVersion APIVersion
}

func NewAPIVersionRouter(apiVersion APIVersion) *APIVersionRouter {
	return &APIVersionRouter{
		mux:        http.NewServeMux(),
		apiVersion: apiVersion,
	}
}

func (r *APIVersionRouter) RegisterRoutes(routes ...Route) {
	for _, route := range routes {
		r.mux.HandleFunc(route.Method+" "+route.Path, route.Handler)
	}
}

func (r *APIVersionRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	serveMux(r.mux, w, req)
}
