package config

// Database type enum.
type DatabaseType int

const (
	Source DatabaseType = iota + 1 // Represents a source database
	Destination // Represents a destination database
)

func (databaseType DatabaseType) String() string {
	return [...]string{"", "Source", "Destination"}[databaseType]
}