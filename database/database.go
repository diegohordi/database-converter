package database

import (
	"database-converter/config"
	"database-converter/errors"
	"database/sql"
)

// Determines a common abstraction for database drivers.
type Database interface {

	// Establishes a database connection using the given configuration.
	Connect(config config.DatabaseConfig) *errors.ApplicationError

	// Disconnects from database.
	Disconnect() *errors.ApplicationError

	// Describes the given source table/collection.
	Describe(source string) (*Table, *errors.ApplicationError)

	// Gets the proper interface for the right database type.
	GetInterface(columnType sql.ColumnType) interface{}

	// Feeds the given channel with rows from the given table.
	GetRows(table Table, columns []string, rowChannel chan interface{})

	// Counts rows from the given table.
	Count(table Table) (int, *errors.ApplicationError)

	// Inserts the row into the given table.
	Insert(table Table, row Row) *errors.ApplicationError
}
