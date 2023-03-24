package nagg

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/inquizarus/nagg/internal/domain"
	"github.com/inquizarus/nagg/pkg/httptools"
)

type Service interface {
	RouteForRequest(*http.Request) (domain.Route, error)
	CreateUpstreamRequest(route domain.Route, request *http.Request) (*http.Request, error)
	GlobalMiddlewares() ([]func(http.Handler) http.Handler, error)
}

type stdService struct {
	config       Config
	routes       []domain.Route
	routesLoaded bool
}

func (s *stdService) loadRoutes() ([]domain.Route, error) {
	if s.routesLoaded {
		return s.routes, nil
	}
	routes, err := s.config.Routes()
	if err != nil {
		s.routesLoaded = false
		return nil, err
	}
	s.routes = routes
	s.routesLoaded = true
	return routes, nil
}

func (s *stdService) RouteForRequest(r *http.Request) (domain.Route, error) {
	routes, err := s.loadRoutes()
	if err != nil {
		return nil, err
	}
	for _, route := range routes {
		if route.Matches(r) {
			return route, nil
		}
	}
	return nil, errors.New("no route found for request")
}

func (s *stdService) CreateUpstreamRequest(route domain.Route, r *http.Request) (*http.Request, error) {

	upstream := route.Address() + r.URL.Path

	if r.URL.Query().Encode() != "" {
		upstream = upstream + "?" + r.URL.Query().Encode()
	}

	upstreamRequest, err := http.NewRequest(r.Method, upstream, r.Body)

	if err != nil {
		return nil, err
	}

	// Copy headers from incoming request to upstream one
	for k, v := range r.Header {
		if len(v) > 0 {
			upstreamRequest.Header.Set(k, v[0])
		}
		if len(v) > 1 {
			// If we have more than one value in a header, plopp 'em values into this request too
			for i := 1; i < len(v); i++ {
				upstreamRequest.Header.Add(k, v[i])
			}
		}
	}

	// Some standard proxy headers
	upstreamRequest.Header.Set("x-forwarded-for", httptools.ClientIP(r))
	upstreamRequest.Header.Set("x-forwarded-host", r.Host)
	upstreamRequest.Header.Set("x-forwarded-proto", r.Proto)
	upstreamRequest.Header.Set("x-request-id", uuid.New().String())

	return upstreamRequest, nil
}

func (s *stdService) GlobalMiddlewares() ([]func(http.Handler) http.Handler, error) {
	return s.config.GlobalMiddlewares()
}

func NewService(config Config) Service {
	return &stdService{
		config:       config,
		routes:       []domain.Route{},
		routesLoaded: false,
	}
}
