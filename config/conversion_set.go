package config

import (
	"database-converter/errors"
	"fmt"
)

// Holds a conversion set configuration.
type ConversionSet struct {
	name        string
	source      string
	destination string
}

// Gets the name of the conversion set.
func (set *ConversionSet) Name() string {
	return set.name
}

// Gets the name of the source table.
func (set *ConversionSet) Source() string {
	return set.source
}

// Gets the name of the destination table.
func (set *ConversionSet) Destination() string {
	return set.name
}

// Validates the database configuration and returns an error just if something goes wrong.
func (set *ConversionSet) validate() *errors.ApplicationError {
	if set.name == "" {
		return errors.WithMessageBuilder("Conversion set name is missing.").Build()
	}
	if set.source == "" {
		return errors.WithMessageBuilder(fmt.Sprintf("Source of conversion set %s is missing.", set.name)).Build()
	}
	if set.destination == "" {
		return errors.WithMessageBuilder(fmt.Sprintf("Destination of conversion set %s is missing.", set.name)).Build()
	}
	return nil
}

