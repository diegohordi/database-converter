package conversion

import (
	"database-converter/config"
	"database-converter/database"
	"database-converter/errors"
	"fmt"
	"os"
	"sync"
	"time"
)

// Service responsible to convert a given conversion set.
type SetService struct {
	sourceConn      database.Database
	destinationConn database.Database
	set             *config.ConversionSet
	statusChannel   chan interface{}
	logFile         *os.File
}

// Builder for this service.
func NewSetService(
	sourceConn database.Database,
	destinationConn database.Database,
	set *config.ConversionSet,
	statusChannel chan interface{}) *SetService {
	return &SetService{
		sourceConn:      sourceConn,
		destinationConn: destinationConn,
		set:             set,
		statusChannel:   statusChannel,
	}
}

// Converts the given conversion set feeding the given channel with the proper status, which means
// an unexpected error if something wrong occurs or a boolean true, if the set was successfully
// converted and false otherwise.
func (service *SetService) Convert() {
	if err := service.createLogFile(); err != nil {
		service.statusChannel <- err
		return
	}
	defer service.closeLogFile()
	service.writeLog(fmt.Sprintf("Start processing of set %s.", service.set.Name()))
	if sourceTable, err := service.sourceConn.Describe(service.set.Source()); err != nil {
		service.writeLog(fmt.Sprintf("An error occurred while describing source table: %s.", err.Error()))
		service.statusChannel <- err
		return
	} else {
		if destinationTable, err := service.destinationConn.Describe(service.set.Destination()); err != nil {
			service.writeLog(fmt.Sprintf("An error occurred while describing destination table: %s.", err.Error()))
			service.statusChannel <- err
			return
		} else {
			if err := service.convertRows(*sourceTable, *destinationTable); err != nil {
				service.writeLog(fmt.Sprintf("An error occurred while processing source table rows: %s", err.Error()))
				service.statusChannel <- err
				return
			}
		}
	}
	service.writeLog(fmt.Sprintf("Finishing conversion of set %s.", service.set.Name()))
	service.statusChannel <- true
}

// Converts the rows from conversion set.
func (service *SetService) convertRows(sourceTable database.Table, destinationTable database.Table) *errors.ApplicationError {
	if totalRows, err := service.sourceConn.Count(sourceTable); err != nil {
		service.writeLog(fmt.Sprintf("An error occurred while couting rows from source table: %s.", err.Error()))
		return err
	} else {
		service.writeLog(fmt.Sprintf("%d row(s) found.", totalRows))
		service.writeLog(fmt.Sprintf("Fetching rows from source."))
		sourceColumns := service.getSelectedSourceColumns(sourceTable)
		var rowChannel = make(chan interface{})
		var rowCounter = 0
		var rowWg sync.WaitGroup
		defer close(rowChannel)
		rowWg.Add(totalRows)
		go service.sourceConn.GetRows(sourceTable, sourceColumns, rowChannel)
		go func() {
			for value := range rowChannel {
				rowWg.Done()
				rowCounter++
				if row, isRow := value.(*database.Row); isRow {
					service.writeLog(fmt.Sprintf("Processing row %d from %d.", rowCounter, totalRows))
					if err := service.destinationConn.Insert(destinationTable, *row); err != nil {
						service.writeLog(fmt.Sprintf("Error inserting row: %s", err.Error()))
					}
				}
				if err, isErr := value.(*errors.ApplicationError); isErr {
					service.writeLog(fmt.Sprintf("An error occurred while processing the row %d: %s.", rowCounter, err.Error()))
				}
			}
		}()
		rowWg.Wait()
	}
	return nil
}

// Gets the selected source columns from the given conversion set. If there are no selected columns
// specified in conversion set, the ones given by the table are returned.
func (service *SetService) getSelectedSourceColumns(table database.Table) []string {
	var columns = make([]string, 0)
	for _, column := range table.Columns() {
		columns = append(columns, column.Name())
	}
	return columns
}

// Creates the log file for the given conversion set.
func (service *SetService) createLogFile() *errors.ApplicationError {
	var err error
	_ = os.Mkdir("logs", 0777)
	service.logFile, err = os.Create(fmt.Sprintf("logs/%s.log", service.set.Name()))
	if err != nil {
		return errors.WithMessageAndSourceErrorBuilder(fmt.Sprintf("Error creating the log file for the set %s.", service.set.Name()), err).Build()
	}
	return nil
}

// Writes a line in the given log file with the given message and current time.
func (service *SetService) writeLog(message string) {
	currentTime := time.Now()
	fmt.Fprintf(service.logFile, "%s: %s\n", currentTime.Format("2006.01.02 15:04:05"), message)
}

// Closes log file.
func (service *SetService) closeLogFile() {
	if service.logFile != nil {
		service.logFile.Close()
	}
}
