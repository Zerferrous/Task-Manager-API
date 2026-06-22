package server

import (
	"fmt"
	"net/http"
)

type ApiVersion string

var (
	ApiVersion1 = ApiVersion("v1")
)

type Router struct {
	*http.ServeMux
	apiVersion ApiVersion
}

func NewRouter(version ApiVersion) *Router {
	return &Router{
		ServeMux:   http.NewServeMux(),
		apiVersion: version,
	}
}

func (r *Router) RegisterRoutes(routes ...Route) {
	for _, route := range routes {
		pattern := fmt.Sprintf("%s %s", route.Method, route.Path)

		r.Handle(pattern, route.Handler)
	}
}
