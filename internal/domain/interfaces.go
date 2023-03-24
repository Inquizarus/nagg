package domain

import (
	"net/http"
)

type Route interface {
	Name() string
	Path() string
	Address() string
	Matches(*http.Request) bool
	Middlewares() ([]func(h http.Handler) http.Handler, error)
	PreMiddlewares() ([]func(h http.Handler) http.Handler, error)
	PostMiddlewares() ([]func(h http.Handler) http.Handler, error)
}

type Predicates interface {
	Path() string
	Matches(*http.Request) bool
}
