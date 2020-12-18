package config

import (
	"database-conversor/errors"
	"fmt"
	"github.com/spf13/viper"
)

type DatabaseConfig struct {
	driver   string
	host     string
	port     int
	user     string
	password string
	database string
}

func (database *DatabaseConfig) GetDriver() string {
	return database.driver
}

func (database *DatabaseConfig) GetHost() string {
	return database.host
}

func (database *DatabaseConfig) GetPort() int {
	return database.port
}

func (database *DatabaseConfig) GetUser() string {
	return database.user
}

func (database *DatabaseConfig) GetPassword() string {
	return database.password
}

func (database *DatabaseConfig) GetDatabase() string {
	return database.database
}

func (database *DatabaseConfig) validate() *errors.ApplicationError {
	if database.driver == "" {
		return errors.BuildApplicationError(nil, "DatabaseConfig driver is missing.", 0)
	}
	if database.host == "" {
		return errors.BuildApplicationError(nil, "DatabaseConfig host is missing.", 0)
	}
	if database.port == 0 {
		return errors.BuildApplicationError(nil, "DatabaseConfig port is missing.", 0)
	}
	if database.user == "" {
		return errors.BuildApplicationError(nil, "DatabaseConfig user is missing.", 0)
	}
	if database.database == "" {
		return errors.BuildApplicationError(nil, "DatabaseConfig name is missing.", 0)
	}
	return nil
}

func parseDatabaseConfig(key string, config *config, viper *viper.Viper) *errors.ApplicationError {
	database := DatabaseConfig{}
	database.driver = viper.GetString(fmt.Sprintf("%s.driver", key))
	database.port = viper.GetInt(fmt.Sprintf("%s.port", key))
	database.host = viper.GetString(fmt.Sprintf("%s.host", key))
	database.user = viper.GetString(fmt.Sprintf("%s.user", key))
	database.password = viper.GetString(fmt.Sprintf("%s.password", key))
	database.database = viper.GetString(fmt.Sprintf("%s.database", key))
	if err := database.validate(); err != nil {
		return err
	}
	switch key {
	case "source":
		config.source = &database
	case "destination":
		config.destination = &database
	}
	return nil
}