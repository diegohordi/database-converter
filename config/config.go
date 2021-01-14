package config

// Holds the application configuration.
type Config struct {
	source      *DatabaseConfig
	destination *DatabaseConfig
	sets        []*ConversionSet
}

// Gets the source database configuration.
func (config *Config) Source() *DatabaseConfig {
	return config.source
}

// Gets the destination database configuration.
func (config *Config) Destination() *DatabaseConfig {
	return config.destination
}

// Gets all conversion sets.
func (config *Config) ConversionSets() []*ConversionSet {
	return config.sets
}
