package database

import (
	"database-converter/config"
	"database-converter/errors"
	"fmt"
	"net/http"
	"strings"
)

/*
Holds the data related to the database connections and contexts.
*/
type Connection struct {
	conn interface{}
	ctx interface{}
	cancel interface{}
}

/*
Represents an abstraction of a database.
*/
type Database interface {

	/*
	Establish a database connection using the given configuration.
	 */
	Connect(config *config.DatabaseConfig) *errors.ApplicationError

	/*
	Disconnect from database.
	 */
	Disconnect() *errors.ApplicationError
}

/*
Database factory, used to create a database representation from the given driver.
*/
func GetDatabase(driver string) (Database, *errors.ApplicationError) {
	switch strings.ToLower(driver) {
	case "mysql":
		return new(MySQL), nil
	case "mongodb":
		return new(MongoDB), nil

	}
	return nil, errors.BuildApplicationError(nil, fmt.Sprintf("The given driver (%s) is not supported yet.", driver), http.StatusBadRequest)
}
