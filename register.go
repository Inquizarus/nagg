package nagg

import (
	"github.com/inquizarus/gosebus"
	"github.com/inquizarus/nagg/pkg/logging"
	"github.com/inquizarus/rwapper/v2"
)

func Register(basePath string, router rwapper.RouterWrapper, eventbus gosebus.Bus, configLoader ConfigLoader, logger logging.Logger) error {

	config, err := configLoader()

	if err != nil {
		return err
	}

	if err := RegisterHTTPHandlers(basePath, router, eventbus, NewService(config), logger); err != nil {
		return err
	}

	return nil
}
