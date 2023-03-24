package nagg

import (
	"io"
	"net/http"

	"github.com/inquizarus/nagg/pkg/logging"
	"github.com/inquizarus/rwapper/v2"
)

// RegisterHTTPHandlers require a router that allows one route to handle all http methods
func RegisterHTTPHandlers(pathPrefix string, router rwapper.RouterWrapper, service Service, logger logging.Logger) {

	handler := makeHandler(service, logger)
	middlewares, err := service.GlobalMiddlewares()

	if err != nil {
		logger.Info("could not load global middlewares for gateway, %s", err.Error())
	}

	router.Handle("", pathPrefix, rwapper.ChainMiddleware(handler, middlewares...))
}

func makeHandler(service Service, logger logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		route, err := service.RouteForRequest(r)

		if err != nil {
			logger.Infof("no route found for ", r.URL.String())
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
