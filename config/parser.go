package config

import (
	"bytes"
	"database-converter/errors"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"strings"
)

// Parser interface.
type Parser interface {
	Parse(registry interface{}) (*Config, *errors.ApplicationError)
}

// Parses and transforms the configuration file into a valid application configuration.
func ParseFile(file string, parser Parser) (*Config, *errors.ApplicationError) {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, errors.WithMessageBuilder(fmt.Sprintf("The file %s is unreachable or doesn't exists.", file)).Build()
	}
	viper.SetConfigType(file[strings.LastIndex(file, ".")+1:])
	if err := viper.ReadConfig(bytes.NewReader(buffer)); err != nil {
		return nil, errors.WithMessageBuilder(fmt.Sprintf("Error loading the config file %s.", file)).Build()
	}
	config, parserErr := parser.Parse(viper.GetViper())
	if parserErr != nil {
		return nil, errors.WithSourceErrorBuilder(parserErr).Build()
	}
	return config, nil
}
