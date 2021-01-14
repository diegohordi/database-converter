package config

import (
	"database-converter/errors"
	"database-converter/utils"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"sync"
)

// Configuration file default parser.
type DefaultParser struct {
	config *Config
}

// Parses the configurations from the given registry.
func (parser *DefaultParser) Parse(registry interface{}) (*Config, *errors.ApplicationError) {
	parser.config = new(Config)
	if viperRegistry, isViper := registry.(*viper.Viper); isViper {
		if err := parser.parseDatabaseConfig(viperRegistry, Source); err != nil {
			return nil, err
		}
		if err := parser.parseDatabaseConfig(viperRegistry, Destination); err != nil {
			return nil, err
		}
		if err := parser.parseConversionsSet(viperRegistry); err != nil {
			return nil, err
		}
		return parser.config, nil
	}
	return nil, errors.WithMessageBuilder("Invalid configuration registry.").Build()
}

// Parses the given database configuration
func (parser *DefaultParser) parseDatabaseConfig(viper *viper.Viper, databaseType DatabaseType) *errors.ApplicationError {
	log.Println(fmt.Sprintf("Parsing %s database configuration.", databaseType))
	database, err := NewDatabaseConfig(
		databaseType,
		GetDatabaseDriver(viper.GetString(fmt.Sprintf("%s.driver", databaseType))),
		viper.GetString(fmt.Sprintf("%s.host", databaseType)),
		viper.GetInt(fmt.Sprintf("%s.port", databaseType)),
		viper.GetString(fmt.Sprintf("%s.user", databaseType)),
		viper.GetString(fmt.Sprintf("%s.password", databaseType)),
		viper.GetString(fmt.Sprintf("%s.database", databaseType)))
	if err != nil {
		return err
	}
	switch databaseType {
	case Source:
		parser.config.source = database
	case Destination:
		parser.config.destination = database
	}
	return nil
}

// Parses all conversion sets
func (parser *DefaultParser) parseConversionsSet(viper *viper.Viper) *errors.ApplicationError {
	log.Println("Parsing conversion sets.")
	sets := viper.GetStringMap("conversion_sets")
	var setsWg sync.WaitGroup
	var countErrors int
	var conversionChannel = make(chan interface{})
	defer close(conversionChannel)
	for setName, definition := range sets {
		setsWg.Add(1)
		go parser.parseConversionSet(setName, definition, conversionChannel)
	}
	go func() {
		for set := range conversionChannel {
			setsWg.Done()
			if conversionSet, isSet := set.(*ConversionSet); isSet {
				parser.config.sets = append(parser.config.sets, conversionSet)
			}
			if err, isError := set.(*errors.ApplicationError); isError {
				log.Println(err.Error())
				countErrors++
			}
		}
	}()
	setsWg.Wait()
	if countErrors > 0 {
		return errors.WithMessageBuilder(fmt.Sprintf("%d set(s) parsed wrong.", countErrors)).Build()
	}
	return nil
}

// Parses the given conversion set and feeds the channel with the parsed conversion set, or an error otherwise.
func (parser *DefaultParser) parseConversionSet(setName string, set interface{}, channel chan interface{}) {
	unboxed, converted := set.(map[string]interface{})
	if !converted {
		channel <- errors.WithMessageBuilder(fmt.Sprintf("Error parsing the conversion set %s.", setName)).Build()
		return
	}
	conversionSet := new(ConversionSet)
	conversionSet.name = setName
	conversionSet.destination = utils.ToString(unboxed["destination"])
	conversionSet.source = utils.ToString(unboxed["source"])
	log.Printf("Parsing conversion set %s.\n", conversionSet.name)
	if err := conversionSet.validate(); err != nil {
		channel <- err
		return
	}
	channel <- conversionSet
}