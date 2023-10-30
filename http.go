package nagg

import (
	"fmt"
	"io"
	"net/http"

	"github.com/inquizarus/gosebus"
	"github.com/inquizarus/nagg/internal/domain"
	"github.com/inquizarus/nagg/internal/events"
	"github.com/inquizarus/nagg/pkg/httptools"
	"github.com/inquizarus/nagg/pkg/logging"
	"github.com/inquizarus/rwapper/v2"
)

// RegisterHTTPHandlers require a router that allows one route to handle all http methods
func RegisterHTTPHandlers(pathPrefix string, router rwapper.RouterWrapper, service Service, logger logging.Logger) error {

	handler := makeHandler(service, logger)
	middlewares, err := service.GlobalMiddlewares()

	if err != nil {
		return fmt.Errorf("could not load global middlewares for gateway, %s", err.Error())
	}

	responseWrapperMiddleware := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rww := httptools.NewResponseWriterWrapper(w, http.StatusOK, logger)
			if h != nil {
				h.ServeHTTP(&rww, r)
			}
		})
	}

	router.Handle("", pathPrefix, rwapper.ChainMiddleware(handler, append(middlewares, responseWrapperMiddleware)...))
	return nil
}

func makeHandler(service Service, logger logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var route domain.Route
		var err error

		defer r.Body.Close()
		defer func() {
			events.PublishRequestWasHandled(gosebus.DefaultBus, w, *r, route)
		}()

		route, err = service.RouteForRequest(r)

		if err != nil {
			logger.Infof("no route found for %s", r.URL.String())
			w.WriteHeader(http.StatusNotFound)
			return
		}

		events.PublishRouteMatchedRequest(gosebus.DefaultBus, route, *r)

		logger.Debugf("found route %s for request with path %s", route.Name(), r.URL.Path)

		preMiddlewares, err := route.PreMiddlewares()

		if err != nil {
			logger.Errorf("could not load pre middlewares, %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		logger.Debug("applying pre middlewares")

		if len(preMiddlewares) > 0 {
			rwapper.ChainMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}), preMiddlewares...).ServeHTTP(w, r)
		}

		if status := w.(httptools.ResponseWriterWrapper).Status(); status != http.StatusOK {
			logger.Infof("response writer status was %d after pre middlewares, returning", status)
			return
		}

		if route.Address() != "" {
			upstreamRequest, err := service.CreateUpstreamRequest(route, r)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			logger.Debugf("performing upstream request to %s", upstreamRequest.URL.String())

			upstreamResponse, err := service.DoRequest(upstreamRequest)

			events.PublishUpstreamRequestWasDone(gosebus.DefaultBus, *upstreamResponse, *upstreamRequest)

			if err != nil {
				logger.Debugf("upstream request resulted in error, %s", err)
				w.WriteHeader(http.StatusNotFound)
				return
			}

			defer upstreamResponse.Body.Close()

			upstreamResponseData, _ := io.ReadAll(upstreamResponse.Body)

			for k, v := range upstreamResponse.Header {
				w.Header().Set(k, v[0])
			}

			w.WriteHeader(upstreamResponse.StatusCode)
			w.Write(upstreamResponseData)
		}

		postMiddlewares, err := route.PostMiddlewares()

		if err != nil {
			logger.Errorf("could not load post middlewares, %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		logger.Debug("applying post middlewares")

		rwapper.ChainMiddleware(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}), postMiddlewares...).ServeHTTP(w, r)

	}
}
