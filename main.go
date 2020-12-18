package main

import (
	"bufio"
	"database-converter/config"
	"database-converter/conversion"
	"database-converter/errors"
	"fmt"
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
	if err := conversion.Prepare(); err != nil {
		err.ThrowPanic()
	}
	fmt.Println("Everything seems to be fine! Do you really want to convert the given databases?")
	fmt.Println("Type yes to continue or anything else to abort.")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		input := scanner.Text()
		if input == "yes" {
			conversion.Start()
		} else {
			fmt.Println("Aborting the conversion. Bye!")
			os.Exit(0)
		}
	}
}
