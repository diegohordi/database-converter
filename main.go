package main

import (
	"database-converter/config"
	"database-converter/conversion"
	"database-converter/errors"
	"net/http"
	"os"
)

/*
Holds the program arguments.
*/
type args struct {
	configFile string // Name of the configuration file.
}

/*
Parse the program arguments, returning the args parsed, otherwise, returns a application error.
 */
func parseArgs(list []string) (*args, *errors.ApplicationError) {
	if len(list) != 2 {
		return nil, errors.BuildApplicationError(nil, "No config file was given.", http.StatusBadRequest)
	}
	args := new(args)
	args.configFile = list[1]
	return args, nil
}

func main() {

	args, err := parseArgs(os.Args)

	if err != nil {
		err.ThrowPanic()
	}

	if err := config.GetInstance().LoadConfig(args.configFile); err != nil {
		err.ThrowPanic()
	}

	if err := conversion.Start(); err != nil {
		err.ThrowPanic()
	}

}
