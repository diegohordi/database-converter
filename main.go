package main

import (
	"bufio"
	"database-converter/config"
	"database-converter/conversion"
	"flag"
	"fmt"
	"os"
)

var (
	configFile = flag.String("config", "", "Configuration file")
)

func main() {
	flag.Parse()
	if *configFile == "" {
		panic("No config file was given.")
	}
	if parsedConfig, err := config.ParseFile(*configFile, &config.DefaultParser{}); err != nil {
		panic(err.Error())
	} else {
		fmt.Println("Everything seems to be fine! Do you really want to convert the given databases?")
		fmt.Println("Type yes to continue or anything else to abort.")
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			input := scanner.Text()
			if input == "yes" {
				service := new(conversion.Service)
				if err := service.Convert(*parsedConfig); err != nil {
					fmt.Println(err.Error())
				}
			} else {
				fmt.Println("Aborting the conversion. Bye!")
				os.Exit(0)
			}
		}
	}

}
