package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/inquizarus/envtools"
	"github.com/inquizarus/gosebus"
	"github.com/inquizarus/nagg"
	"github.com/inquizarus/nagg/pkg/logging"
	"github.com/inquizarus/rwapper/v2/pkg/servemuxwrapper"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
)

func main() {

	var config nagg.Config
	var err error

	logger := logging.NewLogrusLogger(nil, envtools.GetWithFallback(LOG_LEVEL_ENV_KEY, DEFAULT_LOG_LEVEL), nil)

	if config, err = loadGatewayConfig(logger); err != nil {
		logger.Errorf("could not start gateway: %s", err.Error())
		os.Exit(1)
	}

	router := servemuxwrapper.New(nil)
	port := envtools.GetWithFallback(HTTP_PORT_ENV_KEY, DEFAULT_HTTP_PORT)
	basePath := envtools.GetWithFallback(HTTP_BASE_PATH_ENV_KEY, DEFAULT_HTTP_BASE_PATH)

	logger.Debugf("registrering gateway handler on base path %s", basePath)

	router.Handle(http.MethodGet, "/metrics", promhttp.Handler())

	if err = nagg.RegisterHTTPHandlers(basePath, router, nagg.NewService(config), logger); err != nil {
		logger.Errorf("could not start gateway: %s", err.Error())
		os.Exit(1)
	}

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		logger.Infof("gateway starting on port %s", port)

		if err = http.ListenAndServe(":"+port, router); err != nil {
			logger.Error(err)
			signalChannel <- syscall.SIGTERM
		}
	}()

	ctx := context.Background()
	exporter, _ := prometheus.New()
	provider := metric.NewMeterProvider(metric.WithReader(exporter))
	meter := provider.Meter("github.com/inquizarus/nagg/cmd/nagg")
	counter, _ := meter.Int64Counter("nagg_requests_handled_total", metricapi.WithDescription("counter of how many requests that has been handled"))

	gosebus.DefaultBus.Handle(gosebus.NewStandardEventHandler("request_handled", func(e gosebus.Event) {
		counter.Add(ctx, 1)
	}))

	receivedSignal := <-signalChannel

	logger.Infof("received signal '%s', shutting down", receivedSignal)

}

func loadGatewayConfig(logger logging.Logger) (nagg.Config, error) {
	r, found := envtools.GetJSONData(CONFIG_JSON_ENV_KEY)
	if r != nil && found {
		logger.Infof("loading gateway config from JSON in environment variable %s", CONFIG_JSON_ENV_KEY)
		return nagg.JSONConfig(r, nil)
	}
	configPath, err := envtools.GetRequiredOrError(CONFIG_FILE_PATH_ENV_KEY)
	if err != nil {
		return nil, err
	}
	logger.Infof("loading gateway config from file %s", configPath)
	return nagg.JSONConfigFromFile(configPath, nil)
}
