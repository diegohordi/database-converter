package test

import (
	"database-converter/config"
	"database-converter/errors"
	"fmt"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

type ParserTestTable struct {
	name string
	file string
	expected interface{}
}

func Test_DefaultParser(t *testing.T) {
	testTable := make([]ParserTestTable, 0)
	testTable = append(testTable, ParserTestTable{name: "Valid Configuration", file: "valid-configuration.yml", expected: config.Config{}})
	testTable = append(testTable, ParserTestTable{name: "Missing Source", file: "missing-source.yml", expected: errors.ApplicationError{}})
	testTable = append(testTable, ParserTestTable{name: "Missing Source Driver", file: "missing-source-driver.yml", expected: errors.ApplicationError{}})
	testTable = append(testTable, ParserTestTable{name: "Missing Source Host", file: "missing-source-host.yml", expected: errors.ApplicationError{}})
	testTable = append(testTable, ParserTestTable{name: "Missing Source Port", file: "missing-source-port.yml", expected: errors.ApplicationError{}})
	testTable = append(testTable, ParserTestTable{name: "Missing Source User", file: "missing-source-user.yml", expected: errors.ApplicationError{}})
	testTable = append(testTable, ParserTestTable{name: "Missing Source Password", file: "missing-source-password.yml", expected: errors.ApplicationError{}})
	testTable = append(testTable, ParserTestTable{name: "Missing Source Database", file: "missing-source-database.yml", expected: errors.ApplicationError{}})
	testTable = append(testTable, ParserTestTable{name: "Missing Destination", file: "missing-destination.yml", expected: errors.ApplicationError{}})
	testTable = append(testTable, ParserTestTable{name: "Missing Destination Driver", file: "missing-destination-driver.yml", expected: errors.ApplicationError{}})
	testTable = append(testTable, ParserTestTable{name: "Missing Destination Host", file: "missing-destination-host.yml", expected: errors.ApplicationError{}})
	testTable = append(testTable, ParserTestTable{name: "Missing Destination Port", file: "missing-destination-port.yml", expected: errors.ApplicationError{}})
	testTable = append(testTable, ParserTestTable{name: "Missing Destination User", file: "missing-destination-user.yml", expected: errors.ApplicationError{}})
	testTable = append(testTable, ParserTestTable{name: "Missing Destination Password", file: "missing-destination-password.yml", expected: errors.ApplicationError{}})
	testTable = append(testTable, ParserTestTable{name: "Missing Destination Database", file: "missing-destination-database.yml", expected: errors.ApplicationError{}})
	testTable = append(testTable, ParserTestTable{name: "Missing Conversion Set Source", file: "missing-conversion-set-source.yml", expected: errors.ApplicationError{}})
	testTable = append(testTable, ParserTestTable{name: "Missing Conversion Set Destination", file: "missing-conversion-set-destination.yml", expected: errors.ApplicationError{}})
	for _, test := range testTable {
		t.Run( test.name, func(t *testing.T) {
			viper.Reset()
			viper.SetConfigType("yml")
			viper.SetConfigFile(fmt.Sprintf("../test/assets/%s", test.file))
			if err := viper.ReadInConfig(); err != nil {
				t.Error(err)
			}
			parser := &config.DefaultParser{}
			config, err := parser.Parse(viper.GetViper())
			if err != nil {
				assert.IsType(t, test.expected, *err)
			} else {
				assert.IsType(t, test.expected, *config)
			}
		})
	}
}

func Test_DefaultParser_InvalidRegistry(t *testing.T) {
	parser := &config.DefaultParser{}
	config, err := parser.Parse(nil)
	assert.Nil(t, config)
	assert.NotEmpty(t, err)
}
