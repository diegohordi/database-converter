package conversion

import (
	"database-converter/config"
	"database-converter/database"
	"database-converter/errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
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

	// Open a connection with the source database
	sourceConn, err := openDatabaseConnection(config.GetInstance().GetSource())
	defer sourceConn.Disconnect()
	if err != nil {
		return err
	}

	// Open a connection with the destination database
	destinationConn, err := openDatabaseConnection(config.GetInstance().GetDestination())
	defer destinationConn.Disconnect()
	if err != nil {
		return err
	}

	for _, set := range config.GetInstance().GetConversionSets() {
		conversionWg.Add(1)
		go processConversionSet( sourceConn, destinationConn, set, status)
	}

	go func() {
		for processedStatus := range status {
			conversionWg.Done()
			if unboxed, isError := processedStatus.(*errors.ApplicationError); isError {
				log.Printf("%s %s\n", unboxed.GetMessage(), unboxed.GetError().Error())
				countErrors++
			}
			if success, isBoolean := processedStatus.(bool); isBoolean {
				if success {
					countSuccess++
				} else {
					countErrors++
				}
			}
			log.Printf("%d from %d set(s) converted.", countErrors + countSuccess, len(config.GetInstance().GetConversionSets()))
		}
	}()
	conversionWg.Wait()
	log.Printf("Finishing conversion. %d set(s) converted successfully and %d set(s) had errors. Check /logs.", countSuccess, countErrors)
	return nil
}

/*
Process each conversion set, feeding the given channel with the proper status, that means an unexpected error
if something wrong occurs or a boolean true, if the set was successfully converted and false otherwise.
 */
func processConversionSet(sourceConn database.Database, destinationConn database.Database, set *config.ConversionSet, status chan interface{}) {
	logFile, err := createLogFile(set.GetName())
	if err != nil {
		status <- err
		return
	}
	defer logFile.Close()
	writeLog(logFile, fmt.Sprintf("Start processing of set %s.", set.GetName()))
	if sourceTable, err := sourceConn.Describe(set.GetSource()); err != nil {
		writeLog(logFile, fmt.Sprintf("An error occurred while describing source table: %s.", err.GetMessage()))
		status <- err
		return
	} else {
		if destinationTable, err := destinationConn.Describe(set.GetDestination()); err != nil {
			writeLog(logFile, fmt.Sprintf("An error occurred while describing destination table: %s.", err.GetMessage()))
			status <- err
			return
		} else {
			if err := processRows(logFile, sourceConn, destinationConn, set, sourceTable, destinationTable); err != nil {
				writeLog(logFile, fmt.Sprintf("An error occurred while processing source table rows: %s", err.GetMessage()))
				status <- err
				return
			}
		}
	}
	writeLog(logFile, fmt.Sprintf("Finishing conversion of set %s.", set.GetName()))
	status <- true
}

/*
Processing the rows from conversion set.
*/
func processRows(logFile *os.File, sourceConn database.Database, destinationConn database.Database, set *config.ConversionSet, sourceTable *database.Table, destinationTable *database.Table) *errors.ApplicationError{
	if totalRows, err := sourceConn.Count(sourceTable); err != nil {
		writeLog(logFile, fmt.Sprintf("An error occurred while couting rows from source table: %s.", err.GetMessage()))
		return err
	} else {
		writeLog(logFile, fmt.Sprintf("Fetching rows from source."))
		sourceColumns := getSelectedSourceColumns(set, sourceTable)
		var rowChannel = make(chan interface{})
		var rowCounter = 0
		var rowWg sync.WaitGroup
		defer close(rowChannel)
		writeLog(logFile, fmt.Sprintf("%d row(s) found.", totalRows))
		rowWg.Add(totalRows)
		go sourceConn.GetRows(sourceTable, sourceColumns, rowChannel)
		go func() {
			for value := range rowChannel {
				rowWg.Done()
				rowCounter++
				if row, isRow := value.(*database.Row); isRow {
					writeLog(logFile, fmt.Sprintf("Processing row %d from %d.", rowCounter, totalRows))
					if err := destinationConn.Insert(destinationTable, row); err != nil {
						writeLog(logFile, fmt.Sprintf("Error inserting row: %s", err.GetMessage()))
					}
				}
				if err, isErr := value.(*errors.ApplicationError); isErr {
					writeLog(logFile, fmt.Sprintf("An error occurred while processing the row %d: %s.", rowCounter, err.GetMessage()))
				}
			}
		}()
		rowWg.Wait()
	}
	return nil
}

/*
Get the selected source columns from the given conversion set. If there are no selected columns
specified in conversion set, the ones given by the table are returned.
*/
func getSelectedSourceColumns(set *config.ConversionSet, table *database.Table) []string {
	var columns = make([]string, 0)
	for _, column := range table.GetColumns() {
		columns = append(columns, column.GetName())
	}
	return columns
}

/*
Create a log file for the given conversion set.
*/
func createLogFile( setName string ) (*os.File, *errors.ApplicationError) {
	_ = os.Mkdir("logs", 0777)
	logFile, err := os.Create(fmt.Sprintf("logs/%s.log", setName))
	if err != nil {
		return nil, errors.BuildApplicationError(err, fmt.Sprintf("Error creating the log file for the set %s.", setName), http.StatusBadRequest)
	}
	return logFile, nil
}

/*
Writes a line in the given log file with the given message and current time.
*/
func writeLog(logFile *os.File, message string) {
	currentTime := time.Now()
	fmt.Fprintf(logFile, "%s: %s\n", currentTime.Format("2006.01.02 15:04:05"), message)
}

/*
Open a database connection with the given config.DatabaseConfig.
*/
func openDatabaseConnection(config *config.DatabaseConfig) (database.Database, *errors.ApplicationError){
	sourceConnection, err := database.GetDatabase(config.GetDriver())
	if err != nil {
		return nil, err
	}
	sourceConnection.Connect(config)
	return sourceConnection, nil
}