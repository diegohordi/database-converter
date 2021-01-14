package config

import (
	"database-converter/errors"
	"fmt"
)

// Holds the database configurations.
type DatabaseConfig struct {
	databaseType DatabaseType
	driver       DatabaseDriver
	host         string
	port         int
	user         string
	password     string
	name         string
}

// Database config constructor.
func NewDatabaseConfig(
	databaseType DatabaseType,
	driver DatabaseDriver,
	host string,
	port int,
	user string,
	password string,
	name string) (*DatabaseConfig, *errors.ApplicationError) {
	config := &DatabaseConfig{databaseType: databaseType, driver: driver, host: host, port: port, user: user, password: password, name: name}
	if err := config.validate(); err != nil {
		return nil, err
	}
	return config, nil
}

// Gets the database driver.
func (config *DatabaseConfig) Driver() DatabaseDriver {
	return config.driver
}

// Gets the database host.
func (config *DatabaseConfig) Host() string {
	return config.host
}

// Gets the database port.
func (config *DatabaseConfig) Port() int {
	return config.port
}

// Gets the database user.
func (config *DatabaseConfig) User() string {
	return config.user
}

// Gets the user password.
func (config *DatabaseConfig) Password() string {
	return config.password
}

// Gets the database name.
func (config *DatabaseConfig) Database() string {
	return config.name
}

// Validates the database configuration and returns an error just if something goes wrong.
func (config *DatabaseConfig) validate() *errors.ApplicationError {
	if config.databaseType == 0 {
		return errors.WithMessageBuilder(fmt.Sprintf("Database type is missing.")).Build()
	}
	if config.driver == "" {
		return errors.WithMessageBuilder(fmt.Sprintf("%s driver is missing.", config.databaseType)).Build()
	}
	if config.host == "" {
		return errors.WithMessageBuilder(fmt.Sprintf("%s host is missing.", config.databaseType)).Build()
	}
	if config.port == 0 {
		return errors.WithMessageBuilder(fmt.Sprintf("%s port is missing.", config.databaseType)).Build()
	}
	if config.user == "" {
		return errors.WithMessageBuilder(fmt.Sprintf("%s user is missing.", config.databaseType)).Build()
	}
	if config.password == "" {
		return errors.WithMessageBuilder(fmt.Sprintf("%s password is missing.", config.databaseType)).Build()
	}
	if config.name == "" {
		return errors.WithMessageBuilder(fmt.Sprintf("%s name is missing.", config.databaseType)).Build()
	}
	return nil
}
