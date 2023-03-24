package main

import (
	"net/http"
	"os"

	"github.com/inquizarus/envtools"
	"github.com/inquizarus/nagg"
	"github.com/inquizarus/nagg/pkg/logging"
	"github.com/inquizarus/rwapper/v2/pkg/servemuxwrapper"
)

func main() {

	logger := logging.NewLogrusLogger(nil, envtools.GetWithFallback("NAGG_LOG_LEVEL", "info"), nil)

	config, err := loadNaggConfig(logger)

	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}

	router := servemuxwrapper.New(nil)

	if err != nil {
		panic(err)
	}

	port := envtools.GetWithFallback("NAGG_PORT", "8080")
	basePath := envtools.GetWithFallback("NAGG_HTTP_BASE_PATH", "/")

	logger.Debugf("regestering NAGG handler on base path %s", basePath)

	nagg.RegisterHTTPHandlers(basePath, router, nagg.NewService(config), logger)

	logger.Infof("NAGG server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		logger.Error(err)
	}

}

func loadNaggConfig(logger logging.Logger) (nagg.Config, error) {
	r, found := envtools.GetJSONData("NAGG_CONFIG")
	if r != nil && found {
		logger.Info("loading NAGG config from JSON in environment variable NAGG_CONFIG")
		return nagg.JSONConfig(r, nil)
	}
	configPath, err := envtools.GetRequiredOrError("NAGG_CONFIG_PATH")
	if err != nil {
		return nil, err
	}
	logger.Infof("loading NAGG config from file %s", configPath)
	return nagg.JSONConfigFromFile(configPath, nil)
}
