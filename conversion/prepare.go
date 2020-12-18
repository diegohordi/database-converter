package conversion

import (
	"database-converter/config"
	"database-converter/database"
	"database-converter/errors"
	"log"
)

/*
Check the connections to the source and the destination databases.
*/
func testConnections() *errors.ApplicationError {
	if err := testSourceConnection(); err != nil {
		return err
	}
	if err := testDestinationConnection(); err != nil {
		return err
	}
	return nil
}

func testSourceConnection() *errors.ApplicationError{
	db, err := database.GetDatabase(config.GetInstance().GetSource().GetDriver())
	if err != nil {
		return err
	}
	log.Println("Testing connection with the source database.")
	if err := db.Connect(config.GetInstance().GetSource()); err != nil {
		return err
	}
	defer db.Disconnect()
	return nil
}

func testDestinationConnection() *errors.ApplicationError{
	db, err := database.GetDatabase(config.GetInstance().GetDestination().GetDriver())
	if err != nil {
		return err
	}
	log.Println("Testing connection with the destination database.")
	if err := db.Connect(config.GetInstance().GetDestination()); err != nil {
		return err
	}
	defer db.Disconnect()
	return nil
}
