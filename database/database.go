package database

import (
	"database-converter/config"
	"database-converter/errors"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
)

/*
Holds the data related to the database connections and contexts.
*/
type Connection struct {
	config *config.DatabaseConfig
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

	/*
	Describes the given source table/collection.
	 */
	Describe(source string) (*Table, *errors.ApplicationError)

	/*
	Get the proper interface for the right database type.
	 */
	GetInterface(columnType *sql.ColumnType) interface{}

	/*
	Feed the given channel with rows from the given table.
	*/
	GetRows(table *Table, columns []string, rowChannel chan interface{})

	/*
	Count rows from the given table.
	 */
	Count(table *Table) (int, *errors.ApplicationError)

	/*
	Insert the row into the given table.
	 */
	Insert(table *Table, row *Row ) *errors.ApplicationError
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
