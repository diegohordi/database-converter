package config

import (
	"bytes"
	"database-converter/errors"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
)

type config struct {
	source      *DatabaseConfig
	destination *DatabaseConfig
	sets        []*ConversionSet
	sync.RWMutex
}

var (
	instance *config
	once     sync.Once
)

func (config *config) GetSource() *DatabaseConfig {
	return config.source
}

func (config *config) GetDestination() *DatabaseConfig {
	return config.destination
}

func (config *config) GetConversionSets() []*ConversionSet {
	return config.sets
}

func GetInstance() *config {
	once.Do(func() {
		if instance == nil {
			instance = new(config)
		}
	})
	return instance
}

func (config *config) LoadConfig(fileName string) *errors.ApplicationError {
	buffer, err := loadConfigFile(fileName)
	if err != nil {
		return err
	}
	viper.SetConfigType(fileName[strings.LastIndex(fileName, ".")+1:])
	if err := viper.ReadConfig(bytes.NewReader(buffer)); err != nil {
		return errors.BuildApplicationError(err, fmt.Sprintf("Error loading the config file %s.", fileName), http.StatusInternalServerError)
	}
	config.Lock()
	defer config.Unlock()

	log.Println("Parsing source database.")
	if err := parseDatabaseConfig("source", config, viper.GetViper()); err != nil {
		return err
	}

	log.Println("Parsing destination database.")
	if err := parseDatabaseConfig("destination", config, viper.GetViper()); err != nil {
		return err
	}

	log.Println("Parsing conversion sets.")
	sets := viper.GetViper().GetStringMap("conversion_sets")

	var setsWg sync.WaitGroup
	var setsChannel = make(chan interface{})
	defer close(setsChannel)

	for setName, definition := range sets {
		setsWg.Add(1)
		go parseConversionSet(setName, definition, setsChannel)
	}

	go func() {
		for set := range setsChannel {
			setsWg.Done()
			if unboxed, isSet := set.(*ConversionSet); isSet {
				config.sets = append(config.sets, unboxed)
			}
			if unboxed, isError := set.(*errors.ApplicationError); isError {
				log.Println(unboxed.GetMessage())
			}
		}
	}()

	setsWg.Wait()
	return nil
}

func loadConfigFile(file string) ([]byte, *errors.ApplicationError) {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, errors.BuildApplicationError(err, fmt.Sprintf("Error loading the file %s.", file), http.StatusBadRequest)
	}
	return buffer, nil
}
