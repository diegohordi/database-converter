package config

import (
	"database-converter/errors"
	"database-converter/utils"
	"fmt"
	"log"
	"net/http"
)

type ConversionSet struct {
	name        string
	source      string
	destination string
}

func (set *ConversionSet) GetName() string {
	return set.name
}

func (set *ConversionSet) GetSource() string {
	return set.source
}

func (set *ConversionSet) GetDestination() string {
	return set.name
}

func (set *ConversionSet) validate() *errors.ApplicationError {
	if set.name == "" {
		return errors.BuildApplicationError(nil, "Set name is missing.", 0)
	}
	if set.source == "" {
		return errors.BuildApplicationError(nil, fmt.Sprintf("Source of set %s is missing.", set.name), 0)
	}
	if set.destination == "" {
		return errors.BuildApplicationError(nil, fmt.Sprintf("Destination of set %s is missing.", set.name), 0)
	}
	return nil
}

func parseConversionSet(setName string, set interface{}, channel chan interface{}) {
	unboxed, converted := set.(map[string]interface{})
	if !converted {
		channel <- errors.BuildApplicationError(nil, fmt.Sprintf("Error parsing the conversion set %s.", setName), http.StatusBadRequest)
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
