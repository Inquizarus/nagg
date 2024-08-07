package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/inquizarus/envtools"
	"github.com/inquizarus/gosebus"
	"github.com/inquizarus/gosebus/pkg/event"
	ghandler "github.com/inquizarus/gosebus/pkg/handler"
	"github.com/inquizarus/nagg"
	"github.com/inquizarus/nagg/internal/events"
	"github.com/inquizarus/nagg/pkg/httptools"
	"github.com/inquizarus/nagg/pkg/logging"
	"github.com/inquizarus/rwapper/v2"
	"github.com/inquizarus/rwapper/v2/pkg/servemuxwrapper"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/prometheus"
	metricapi "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
)

const (
	DEFAULT_HTTP_PORT      string = "8080"
	DEFAULT_HTTP_BASE_PATH string = "/"
	DEFAULT_LOG_LEVEL      string = "info"

	// Different key strings to retrieve information from the runtime
	// environment
	LOG_LEVEL_ENV_KEY        string = "NAGG_LOG_LEVEL"
	HTTP_PORT_ENV_KEY        string = "NAGG_HTTP_PORT"
	HTTP_BASE_PATH_ENV_KEY   string = "NAGG_HTTP_BASE_PATH"
	CONFIG_JSON_ENV_KEY      string = "NAGG_CONFIG_JSON"
	CONFIG_FILE_PATH_ENV_KEY string = "NAGG_CONFIG_FILE_PATH"

	ADD_ENDPOINT_TO_METRICS_ENV_KEY string = "NAGG_ADD_ENDPOINT_TO_METRICS"
)

func main() {

	//	var config nagg.Config
	var err error

	logger := logging.NewLogrusLogger(nil, envtools.GetWithFallback(LOG_LEVEL_ENV_KEY, DEFAULT_LOG_LEVEL), nil)
	router := servemuxwrapper.New(nil)
	port := envtools.GetWithFallback(HTTP_PORT_ENV_KEY, DEFAULT_HTTP_PORT)
	eventbus := gosebus.New()

	eventbus.On("nagg_*", func(e event.Event) {
		logger.Debugf("NAGG event %s", e.Name())
	})

	/*
	*	Metrics configuration
	 */
	ctx := context.Background()
	enableMetrics(ctx, router, eventbus)

	/*
	 * Nagg handler configuration
	 */

	nagg.Register(
		envtools.GetWithFallback(HTTP_BASE_PATH_ENV_KEY, DEFAULT_HTTP_BASE_PATH),
		router,
		eventbus,
		makeLoadGatewayConfig(logger),
		logger,
	)

	/*
	* Start server and handle graceful shutdowns
	 */

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		logger.Infof("gateway starting on port %s", port)

		if err = http.ListenAndServe(":"+port, router); err != nil {
			logger.Error(err)
			signalChannel <- syscall.SIGTERM
		}
	}()

	receivedSignal := <-signalChannel

	logger.Infof("received signal '%s', shutting down", receivedSignal)

}

func enableMetrics(ctx context.Context, router rwapper.RouterWrapper, eventbus gosebus.Bus) {
	addEndpointToMetrics := envtools.GetWithFallback(ADD_ENDPOINT_TO_METRICS_ENV_KEY, "")

	router.Handle(http.MethodGet, "/metrics", promhttp.Handler())

	exporter, _ := prometheus.New()
	provider := metric.NewMeterProvider(metric.WithReader(exporter))
	meter := provider.Meter("github.com/inquizarus/nagg/cmd/nagg")

	requestCounter, _ := meter.Int64Counter("nagg_request_total", metricapi.WithDescription("counter of how many requests that has been handled"))
	httpStatusCounter, _ := meter.Int64Counter("nagg_http_status", metricapi.WithDescription("counter of http statuses"))
	routeCounter, _ := meter.Int64Counter("nagg_route_hits", metricapi.WithDescription("counter of route hits"))

	eventbus.Handle(ghandler.New(func(e event.Event) {

		data := e.Data().(events.RequestHandled)

		requestCounter.Add(ctx, 1)

		if addEndpointToMetrics != "" {
			httpStatusCounter.Add(ctx, 1, metricapi.WithAttributes(
				attribute.Int("status", data.ResponseWriter.(httptools.ResponseWriterWrapper).Status()),
				attribute.String("uri", data.Request.RequestURI),
			))
		} else {
			httpStatusCounter.Add(ctx, 1, metricapi.WithAttributes(
				attribute.Int("status", data.ResponseWriter.(httptools.ResponseWriterWrapper).Status()),
			))
		}

		if data.Route != nil {
			routeCounter.Add(ctx, 1, metricapi.WithAttributes(
				attribute.String("name", data.Route.Name()),
				attribute.String("path", data.Route.Path()),
			))
		}

	},
		ghandler.WithPattern(events.REQUEST_WAS_HANDLED_EVENT),
	))
}

func makeLoadGatewayConfig(logger logging.Logger) nagg.ConfigLoader {
	return func() (nagg.Config, error) {
		r, found := envtools.GetJSONData(CONFIG_JSON_ENV_KEY)
		if r != nil && found {
			logger.Infof("parsing config from json string in environment variable %s", CONFIG_JSON_ENV_KEY)
			return nagg.JSONConfig(r, nil)
		}
		configPath, err := envtools.GetRequiredOrError(CONFIG_FILE_PATH_ENV_KEY)
		if err != nil {
			return nil, err
		}
		logger.Infof("loading nagg config json file with path from environment variable %s", CONFIG_FILE_PATH_ENV_KEY)
		return nagg.JSONConfigFromFile(configPath, nil)
	}
}
