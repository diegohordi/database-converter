package conversion

import (
	"database-converter/config"
	"database-converter/errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
)

/*
Prepare the conversion.
*/
func Prepare() *errors.ApplicationError {
	if err := testConnections(); err != nil {
		return err
	}
	return nil
}

/*
Start the conversion.
*/
func Start() *errors.ApplicationError {
	var conversionWg sync.WaitGroup
	var status = make(chan interface{})
	log.Println("Starting conversion.")
	var countErrors int
	var countSuccess int
	for _, set := range config.GetInstance().GetConversionSets() {
		conversionWg.Add(1)
		go processConversionSet(set, status)
	}
	go func() {
		for processedStatus := range status {
			conversionWg.Done()
			if unboxed, isError := processedStatus.(*errors.ApplicationError); isError {
				log.Printf("%s %s\n", unboxed.GetMessage(), unboxed.GetError().Error())
				countErrors++
			} else {
				countSuccess++
			}
			log.Printf("%d from %d set(s) converted.", countErrors + countSuccess, len(config.GetInstance().GetConversionSets()))
		}
	}()
	conversionWg.Wait()
	log.Printf("Finishing conversion. %d set(s) converted successfully and %d set(s) had errors. Check /logs.", countSuccess, countErrors)
	return nil
}

/*
Process each conversion set, feeding the given channel with the proper status, that means an error
if something wrong occurs or a boolean true, if the set was successfully converted.
 */
func processConversionSet(set *config.ConversionSet, status chan interface{}) {
	_ = os.Mkdir("logs", 0777)
	logFile, err := os.Create(fmt.Sprintf("logs/%s.log", set.GetName()))
	defer logFile.Close()
	if err != nil {
		status <- errors.BuildApplicationError(err, fmt.Sprintf("Error creating the log file for the set %s.", set.GetName()), http.StatusBadRequest)
		return
	}
	
	status <- true
}
