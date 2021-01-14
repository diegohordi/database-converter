package conversion

import (
	"database-converter/config"
	"database-converter/database"
	"database-converter/errors"
	"log"
	"sync"
)

// Service responsible to convert the databases using the given configuration.
type Service struct {
	config config.Config
	sourceConn database.Database
	destinationConn database.Database
}

// Converts the databases.
func (service *Service) Convert(config config.Config) *errors.ApplicationError {
	log.Println("Starting conversion.")
	service.config = config
	if err := service.establishConnections(); err != nil {
		return err
	} else {
		defer service.closeConnections()
		if err := service.convertConversionSets(); err != nil {
			return err
		}
	}
	return nil
}

// Close connections with both, source and destination databases
func (service *Service) closeConnections() {
	if service.sourceConn != nil {
		if err := service.sourceConn.Disconnect(); err != nil {
			log.Println(err.Error())
		}
	}
	if service.destinationConn != nil {
		if err := service.destinationConn.Disconnect(); err != nil {
			log.Println(err.Error())
		}
	}
}

// Establishes connections with both, source and destination databases.
func (service *Service) establishConnections() *errors.ApplicationError {
	log.Println("Establishing database connections.")
	var err *errors.ApplicationError
	if service.sourceConn, err = service.openDatabaseConnection(*service.config.Source()); err != nil {
		return err
	}
	if service.destinationConn, err = service.openDatabaseConnection(*service.config.Destination()); err != nil {
		return err
	}
	return nil
}

// Opens a database connection according to the given type.
func (service *Service) openDatabaseConnection(databaseConfig config.DatabaseConfig) (database.Database, *errors.ApplicationError) {
	if db, err := database.GetDatabase(databaseConfig.Driver()); err != nil {
		return nil, errors.WithSourceErrorBuilder(err).Build()
	} else {
		if err := db.Connect(databaseConfig); err != nil {
			return nil, errors.WithSourceErrorBuilder(err).Build()
		} else {
			return db, nil
		}
	}
}

// Converts all conversion sets configured.
func (service *Service) convertConversionSets() *errors.ApplicationError{
	var conversionWg sync.WaitGroup
	var status = make(chan interface{})
	var countErrors int
	var countSuccess int

	for _, set := range service.config.ConversionSets() {
		conversionWg.Add(1)
		setService := NewSetService(service.sourceConn, service.destinationConn, set, status)
		go setService.Convert()
	}

	go func() {
		for processedStatus := range status {
			conversionWg.Done()
			if unboxed, isError := processedStatus.(*errors.ApplicationError); isError {
				log.Printf("%s %s\n", unboxed.Message(), unboxed.SourceError().Error())
				countErrors++
			}
			if success, isBoolean := processedStatus.(bool); isBoolean {
				if success {
					countSuccess++
				} else {
					countErrors++
				}
			}
			log.Printf("%d from %d set(s) converted.", countErrors + countSuccess, len(service.config.ConversionSets()))
		}
	}()

	conversionWg.Wait()
	log.Printf("%d set(s) converted successfully and %d set(s) had errors. Check /logs directory.", countSuccess, countErrors)
	return nil
}