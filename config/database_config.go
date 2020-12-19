package config

import (
	"database-converter/errors"
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

func (database *DatabaseConfig) validate(databaseType string) *errors.ApplicationError {
	if database.driver == "" {
		return errors.BuildApplicationError(nil, fmt.Sprintf("%s driver is missing.", databaseType), 0)
	}
	if database.host == "" {
		return errors.BuildApplicationError(nil, fmt.Sprintf("%s host is missing.", databaseType), 0)
	}
	if database.port == 0 {
		return errors.BuildApplicationError(nil, fmt.Sprintf("%s port is missing.", databaseType), 0)
	}
	if database.user == "" {
		return errors.BuildApplicationError(nil, fmt.Sprintf("%s user is missing.", databaseType), 0)
	}
	if database.password == "" {
		return errors.BuildApplicationError(nil, fmt.Sprintf("%s password is missing.", databaseType), 0)
	}
	if database.database == "" {
		return errors.BuildApplicationError(nil, fmt.Sprintf("%s database is missing.", databaseType), 0)
	}
	return nil
}

func parseDatabaseConfig(databaseType string, config *config, viper *viper.Viper) *errors.ApplicationError {
	database := DatabaseConfig{}
	database.driver = viper.GetString(fmt.Sprintf("%s.driver", databaseType))
	database.port = viper.GetInt(fmt.Sprintf("%s.port", databaseType))
	database.host = viper.GetString(fmt.Sprintf("%s.host", databaseType))
	database.user = viper.GetString(fmt.Sprintf("%s.user", databaseType))
	database.password = viper.GetString(fmt.Sprintf("%s.password", databaseType))
	database.database = viper.GetString(fmt.Sprintf("%s.database", databaseType))
	if err := database.validate(databaseType); err != nil {
		return err
	}
	switch databaseType {
	case "source":
		config.source = &database
	case "destination":
		config.destination = &database
	}
	return nil
}