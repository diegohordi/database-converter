package database

import "database-converter/config"

// Holds the data related to the database connections and contexts.
type Connection struct {
	config config.DatabaseConfig
	conn interface{}
	ctx interface{}
	cancel interface{}
}
