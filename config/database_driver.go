package config

import "strings"

// Database driver enum.
type DatabaseDriver string

const (
	MySQL DatabaseDriver = "MySQL"
	MongoDB = "MongoDB"
)

// Gets the database driver by the given driver name. Returns an empty string if no equivalent has been found.
func GetDatabaseDriver(driverName string) DatabaseDriver{
	driverName = strings.ToLower(driverName)
	switch driverName {
	case "mysql":
		return MySQL
	case "mongodb":
		return MongoDB
	default:
		return ""
	}
}