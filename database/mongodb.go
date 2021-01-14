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
)

type MongoDBImpl struct {
	connection *Connection
}

func (db *MongoDBImpl) Connect(config config.DatabaseConfig) *errors.ApplicationError {
	serverUrl := fmt.Sprintf("mongodb://%s:%d", config.Host(), config.Port())
	credential := options.Credential{
		Username: config.User(),
		Password: config.Password(),
	}
	clientOptions := options.Client().ApplyURI(serverUrl).SetAuth(credential)
	ctx := context.Background()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return errors.WithMessageAndSourceErrorBuilder("Error creating the MongoDB client.", err).Build()
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return errors.WithMessageAndSourceErrorBuilder("Error establishing connection with database.", err).Build()
	}
	db.connection = &Connection{config: config, conn: client, ctx: ctx}
	return nil
}

func (db *MongoDBImpl) Disconnect() *errors.ApplicationError {
	if instance, check := db.connection.conn.(*mongo.Client); check {
		if context, check := db.connection.ctx.(context.Context); check {
			if err := instance.Disconnect(context); err != nil {
				return errors.WithMessageAndSourceErrorBuilder("Error closing connection.", err).Build()
			}
		}
	}
	return nil
}

func (db *MongoDBImpl) Describe(source string) (*Table, *errors.ApplicationError) {
	return &Table{name: source}, nil
}

func (db *MongoDBImpl) GetInterface(columnType sql.ColumnType) interface{} {
	panic("implement me")
}

func (db *MongoDBImpl) GetRows(table Table, columns []string, rowChannel chan interface{}) {
	panic("implement me")
}

func (db *MongoDBImpl) Count(table Table) (int, *errors.ApplicationError) {
	panic("implement me")
}

func (db *MongoDBImpl) Insert(table Table, row Row) *errors.ApplicationError {
	if instance, check := db.connection.conn.(*mongo.Client); check {
		collection := instance.Database(db.connection.config.Database()).Collection(table.name)
		if _, err := collection.InsertOne(context.TODO(), row.data); err != nil {
			return errors.WithMessageAndSourceErrorBuilder("An error occurred while inserting row.", err).Build()
		}
	}
	return nil
}
