package nagg

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/inquizarus/nagg/internal/domain"
	"github.com/inquizarus/nagg/pkg/middlewares"
)

type jsonMiddleware struct {
	Name  string        `json:"name"`
	Phase string        `json:"phase"`
	Args  []interface{} `json:"args"`
}

type jsonPredicates struct {
	Path string `json:"path"`
}

type jsonRoute struct {
	Name        string           `json:"name"`
	Address     string           `json:"address"`
	Predicates  jsonPredicates   `json:"predicates"`
	Middlewares []jsonMiddleware `json:"middlewares"`
}

type jsonHTTP struct {
	DisableTLSVerification bool `json:"disable_tls_verification"`
	Timeout                int  `json:"timeout"`
}

type jsonMetrics struct {
	Enabled              bool   `json:"enabled"`
	AddEndpointToMetrics bool   `json:"add_endpoint_to_metrics"`
	Path                 string `json:"path"`
}

type jsonGateway struct {
	HTTP       jsonHTTP         `json:"http"`
	Routes     []jsonRoute      `json:"routes"`
	Middleware []jsonMiddleware `json:"middlewares"`
}

type jsonConfig struct {
	middlewareLoader MiddlewareLoader
	Gateway          jsonGateway `json:"gateway"`
}

type ConfigLoader func() (Config, error)

func (config *jsonConfig) Routes() ([]domain.Route, error) {
	routes := []domain.Route{}
	for _, route := range config.Gateway.Routes {
		middlewares := config.middlewaresFromConfigs(route.Middlewares, config.middlewareLoader)
		routes = append(routes, domain.NewRoute(route.Name, route.Address, domain.NewPredicates(route.Predicates.Path), middlewares))
	}
	return routes, nil
}

func (config *jsonConfig) GlobalMiddlewares() ([]func(http.Handler) http.Handler, error) {
	middlewares := config.middlewaresFromConfigs(config.Gateway.Middleware, config.middlewareLoader)
	return append(middlewares["pre"], middlewares["post"]...), nil
}

func (config *jsonConfig) HTTPClient() (*http.Client, error) {

	transport := http.DefaultTransport.(*http.Transport).Clone()

	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: config.Gateway.HTTP.DisableTLSVerification}

	client := &http.Client{
		Transport: transport,
		Timeout:   time.Duration(config.Gateway.HTTP.Timeout) * time.Second,
	}

	return client, nil
}

func (cfg *jsonConfig) middlewaresFromConfigs(configs []jsonMiddleware, middlewareLoader MiddlewareLoader) map[string][]func(http.Handler) http.Handler {
	preMiddlewares := []func(http.Handler) http.Handler{}
	postMiddlewares := []func(http.Handler) http.Handler{}
	for _, config := range configs {
		var middleware func(http.Handler) http.Handler

		if middlewareLoader != nil {
			middleware = middlewareLoader(config.Name, config.Args)
		}

		if middleware == nil {
			middleware = middlewares.DefaultLoader(config.Name, config.Args)
		}

		if middleware != nil {
			if config.Phase == "post" {
				postMiddlewares = append(postMiddlewares, middleware)
				continue
			}
			preMiddlewares = append(preMiddlewares, middleware)
		}
	}
	return map[string][]func(http.Handler) http.Handler{
		"pre":  preMiddlewares,
		"post": postMiddlewares,
	}
}

func JSONConfig(r io.Reader, middlewareLoader MiddlewareLoader) (Config, error) {
	if r == nil {
		return nil, errors.New("reader was nil, config could not be loaded")
	}
	config := jsonConfig{
		middlewareLoader: middlewareLoader,
	}
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &config)
	return &config, err
}

func JSONConfigFromFile(path string, middlewareLoader MiddlewareLoader) (Config, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(file)

	if err != nil {
		return nil, err
	}

	file.Close()

	return JSONConfig(bytes.NewReader(data), middlewareLoader)
}

func MakeJSONConfigLoader(r io.Reader, middlewareLoader MiddlewareLoader) ConfigLoader {
	return func() (Config, error) {
		return JSONConfig(r, middlewareLoader)
	}
}

func MakeJSONConfigFromFileLoader(path string, middlewareLoader MiddlewareLoader) ConfigLoader {
	return func() (Config, error) {
		return JSONConfigFromFile(path, middlewareLoader)
	}
}
