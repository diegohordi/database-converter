package database

import (
	"database-converter/config"
	"database-converter/errors"
	"fmt"
)

// Database factory, used to get the database implementation of the given driver.
func GetDatabase(driver config.DatabaseDriver) (Database, *errors.ApplicationError) {
	switch driver {
	case config.MySQL:
		return new(MySQLImpl), nil
	case config.MongoDB:
		return new(MongoDBImpl), nil
	}
	return nil, errors.WithMessageBuilder(fmt.Sprintf("The given driver (%s) is not supported yet.", driver)).Build()
}
