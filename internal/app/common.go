package app

import (
	database_postgres "github.com/markgredasov/viewcode-sandbox/internal/database/postgres"
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

func (a *App) initPostgresPool() error {
	databaseConfig := database_postgres.NewConfigMust()
	postgresPool, err := database_postgres.NewConnectionPool(a.ctx, databaseConfig)
	if err != nil {
		return err
	}

	a.pool = postgresPool

	return nil
}
