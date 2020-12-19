package database

import (
	"context"
	"database-converter/config"
	"database-converter/errors"
	"database/sql"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"net/http"
)

type MongoDB struct {
	connection *Connection
}

func (db *MongoDB) Connect(config *config.DatabaseConfig) *errors.ApplicationError {
	serverUrl := fmt.Sprintf("mongodb://%s:%d", config.GetHost(), config.GetPort())
	credential := options.Credential{
		Username: config.GetUser(),
		Password: config.GetPassword(),
	}
	clientOptions := options.Client().ApplyURI(serverUrl).SetAuth(credential)
	ctx := context.Background()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return errors.BuildApplicationError(err, "Error creating the MongoDB client.", http.StatusBadRequest)
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return errors.BuildApplicationError(err, "Error establishing connection with database.", http.StatusBadRequest)
	}
	db.connection = &Connection{config: config, conn: client, ctx: ctx}
	return nil
}

func (db *MongoDB) Disconnect() *errors.ApplicationError{
	if instance, check := db.connection.conn.(*mongo.Client); check {
		if context, check := db.connection.ctx.(context.Context); check {
			if err:= instance.Disconnect(context); err != nil {
				return errors.BuildApplicationError(err, "Error closing connection.", http.StatusInternalServerError)
			}
		}
	}
	return nil
}

func (db *MongoDB) Describe(source string) (*Table, *errors.ApplicationError) {
	return &Table{name: source}, nil
}

func (db *MongoDB) GetInterface(columnType *sql.ColumnType) interface{} {
	panic("implement me")
}

func (db *MongoDB) GetRows(table *Table, columns []string, rowChannel chan interface{}){
	panic("implement me")
}

func (db *MongoDB) Count(table *Table) (int, *errors.ApplicationError) {
	panic("implement me")
}

func (db *MongoDB) Insert(table *Table, row *Row) *errors.ApplicationError {
	if instance, check := db.connection.conn.(*mongo.Client); check {
		collection := instance.Database(db.connection.config.GetDatabase()).Collection(table.name)
		if _, err := collection.InsertOne(context.TODO(), row.data); err != nil {
			return errors.BuildApplicationError(err, "An error occurred while inserting row.", http.StatusInternalServerError)
		}
	}
	return nil
}

