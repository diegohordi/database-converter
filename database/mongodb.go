package database

import (
	"context"
	"database-converter/config"
	"database-converter/errors"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"net/http"
	"time"
)

type MongoDB struct {
	connection *Connection
}

func (db *MongoDB) Connect(config *config.DatabaseConfig) *errors.ApplicationError {
	serverUrl := fmt.Sprintf("mongodb://%s:%d", config.GetHost(), config.GetPort())
	clientOptions := options.Client().ApplyURI(serverUrl)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return errors.BuildApplicationError(err, "Error creating the MongoDB client.", http.StatusBadRequest)
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return errors.BuildApplicationError(err, "Error establishing connection with database.", http.StatusBadRequest)
	}
	db.connection = &Connection{conn: client, ctx: ctx, cancel: cancel}
	return nil
}

func (db *MongoDB) Disconnect() *errors.ApplicationError{
	if instance, check := db.connection.conn.(mongo.Client); check {
		if context, check := db.connection.ctx.(context.Context); check {
			if err:= instance.Disconnect(context); err != nil {
				return errors.BuildApplicationError(err, "Error closing connection.", http.StatusInternalServerError)
			}
		}
	}
	return nil
}
