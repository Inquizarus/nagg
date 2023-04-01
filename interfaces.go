package nagg

import (
	"net/http"

	"github.com/inquizarus/nagg/internal/domain"
)

// MiddlewareLoader must return a middleware handler based on the given name
// with passed arguments correctly injected and type-casted for that specific middleware
type MiddlewareLoader func(name string, args ...interface{}) func(http.Handler) http.Handler

// RoutesLoader retrieves routes from some location and tries to use the passed middleware loader
// for each route, if the passed middleware loader is nil or results in a nil middleware handler
// the default middleware loader will be used
type RoutesLoader func(middlewareLoader MiddlewareLoader) ([]domain.Route, error)

type Config interface {
	Routes() ([]domain.Route, error)
	GlobalMiddlewares() ([]func(http.Handler) http.Handler, error)
	HTTPClient() (*http.Client, error)
}
