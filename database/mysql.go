package database

import (
	"database-converter/config"
	"database-converter/errors"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)

type MySQL struct {
	connection *Connection
}

func (db *MySQL) Connect(config *config.DatabaseConfig) *errors.ApplicationError {
	connString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", config.GetUser(), config.GetPassword(), config.GetHost(), config.GetPort(), config.GetDatabase())
	conn, err := sql.Open("mysql", connString)
	if err != nil {
		return errors.BuildApplicationError(err, "Error establishing connection with database.", http.StatusBadRequest)
	}
	if err:= conn.Ping(); err != nil {
		return errors.BuildApplicationError(err, "The database is unreachable.", http.StatusBadRequest)
	}
	db.connection = &Connection{conn: conn}
	return nil
}

func (db *MySQL) Disconnect() *errors.ApplicationError{
	if instance, check := db.connection.conn.(sql.DB); check {
		if err:= instance.Close(); err != nil {
			return errors.BuildApplicationError(err, "Error closing connection.", http.StatusInternalServerError)
		}
	}
	return nil
}


