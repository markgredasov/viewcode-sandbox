package app

import (
	"github.com/markgredasov/viewcode-sandbox/internal/logger"
)

func (a *App) initLogger() error {
	loggerConfig := logger.NewConfigMust()
	logger, err := logger.NewLogger(loggerConfig)
	if err != nil {
		return err
	}
	a.logger = logger

	return nil
}
