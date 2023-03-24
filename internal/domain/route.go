package domain

import (
	"net/http"
	"os"
	"strings"
)

type route struct {
	name        string
	address     string
	predicates  Predicates
	middlewares map[string][]func(http.Handler) http.Handler
}

func (r *route) Name() string {
	return r.name
}

func (r *route) Path() string {
	return r.predicates.Path()
}

func (r *route) Address() string {
	if strings.HasPrefix(r.address, "env:") {
		return os.Getenv(strings.TrimPrefix(r.address, "env:"))
	}
	return r.address
}

func (r *route) Matches(req *http.Request) bool {
	return r.predicates.Matches(req)
}

func (r *route) Middlewares() ([]func(h http.Handler) http.Handler, error) {
	return r.middlewaresByPhase("")
}

func (r *route) PreMiddlewares() ([]func(h http.Handler) http.Handler, error) {
	return r.middlewaresByPhase("pre")
}

func (r *route) PostMiddlewares() ([]func(h http.Handler) http.Handler, error) {
	return r.middlewaresByPhase("post")
}

func (r *route) middlewaresByPhase(phase string) ([]func(h http.Handler) http.Handler, error) {
	if middlewares, ok := r.middlewares[phase]; ok {
		return middlewares, nil
	}

	middlewares := []func(h http.Handler) http.Handler{}

	for _, mws := range r.middlewares {
		middlewares = append(middlewares, mws...)
	}

	return middlewares, nil
}

func NewRoute(name, address string, predicates Predicates, middlewares map[string][]func(http.Handler) http.Handler) Route {
	return &route{
		name:        name,
		predicates:  predicates,
		address:     address,
		middlewares: middlewares,
	}
}
