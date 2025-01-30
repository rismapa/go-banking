package config

import (
	"github.com/okyws/go-banking/domain"

	"github.com/jmoiron/sqlx"
	logger "github.com/okyws/go-banking-lib/config"
)

/*
 * Implemtasi database dengan config dari .yaml
 */
func NewDBConnectionYAML() (*sqlx.DB, error) {
	config, err := domain.GetConfig()
	if err != nil {
		logger.GetLog().Fatal().Err(err).Msg("Failed to get config")
	}

	db, err := sqlx.Connect("mysql", config.GetDatabaseConfig())
	if err != nil {
		logger.GetLog().Fatal().Err(err).Msg("Failed to connect database")
	} else {
		logger.GetLog().Info().Msg("Database connected")
	}

	return db, nil
}

/*
 * Use database config from .env
 */
func NewDBConnectionENV() (*sqlx.DB, error) {
	config, err := domain.GetConfig()
	if err != nil {
		logger.GetLog().Fatal().Err(err).Msg("Failed to get config")
	}

	db, err := sqlx.Connect("mysql", config.GetDatabaseENVConfig())
	if err != nil {
		logger.GetLog().Fatal().Err(err).Msg("Failed to connect database")
	} else {
		logger.GetLog().Info().Msg("Database connected")
	}

	return db, nil
}
