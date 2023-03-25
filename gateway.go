package nagg

import (
	"io"
	"net/http"

	"github.com/inquizarus/nagg/pkg/logging"
	"github.com/inquizarus/rwapper/v2"
)

type responseWriterWrapper struct {
	http.ResponseWriter
	StatusCode int
	log        logging.Logger
}

func (rww *responseWriterWrapper) WriteHeader(code int) {
	rww.log.Debugf("changing response status code from %d to %d", rww.StatusCode, code)
	rww.StatusCode = code
	rww.ResponseWriter.WriteHeader(code)
}

// RegisterHTTPHandlers require a router that allows one route to handle all http methods
func RegisterHTTPHandlers(pathPrefix string, router rwapper.RouterWrapper, service Service, logger logging.Logger) {

	handler := makeHandler(service, logger)
	middlewares, err := service.GlobalMiddlewares()

	responseWrapperMiddleware := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rww := responseWriterWrapper{w, http.StatusOK, logger}
			if h != nil {
				h.ServeHTTP(&rww, r)
			}
		})
	}

	if err != nil {
		logger.Info("could not load global middlewares for gateway, %s", err.Error())
	}

	router.Handle("", pathPrefix, rwapper.ChainMiddleware(handler, append(middlewares, responseWrapperMiddleware)...))
}

func makeHandler(service Service, logger logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		route, err := service.RouteForRequest(r)

		if err != nil {
			logger.Infof("no route found for %s", r.URL.String())
			w.WriteHeader(http.StatusNotFound)
			return
		}

		preMiddlewares, err := route.PreMiddlewares()

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Apply route specific preMiddlewares,
		rwapper.ChainMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}), preMiddlewares...).ServeHTTP(w, r)

		if status := w.(*responseWriterWrapper).StatusCode; status != http.StatusOK {
			logger.Infof("response writer status was %d after pre middlewares, returning", status)
			return
		}

		upstreamRequest, err := service.CreateUpstreamRequest(route, r)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		logging.DefaultLogger.Debugf("performing upstream request to %s", upstreamRequest.URL.String())

		upstreamResponse, err := http.DefaultClient.Do(upstreamRequest)

		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		defer upstreamResponse.Body.Close()

		postMiddlewares, err := route.PostMiddlewares()

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		upstreamResponseData, _ := io.ReadAll(upstreamResponse.Body)

		for k, v := range upstreamResponse.Header {
			w.Header().Set(k, v[0])
		}

		w.WriteHeader(upstreamResponse.StatusCode)
		w.Write(upstreamResponseData)

		rwapper.ChainMiddleware(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}), postMiddlewares...).ServeHTTP(w, r)

	}
}
