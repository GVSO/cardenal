package db

import (
	"context"
	"fmt"

	"github.com/gvso/cardenal/src/app/settings"
	"github.com/mongodb/mongo-go-driver/mongo"
)

var connected = false

// Client is the client connection to MongoDB
var Client MongoClient

// Database is the connection to the app database.
var Database MongoDatabase

/**
 * Helper functions from external packages.
 *
 * For testing purposes, variables are declared which are defined as pointers to
 * function. The advantage of doing this is that these variables can later be
 * overwritten in the testing files.
 */
var newMongoClient = mongo.NewClient

// StartConnection starts the connection to database.
var StartConnection = func() error {

	if !connected {
		var err error

		Client, err = getMongoClient()
		if err != nil {
			return err
		}

		err = Client.Connect(context.TODO())

		if err != nil {
			return err
		}

		Database = Client.Database(settings.MongoDB.Database)

		connected = true
	}

	return nil
}

// GetCollection gets a collection in the database.
var GetCollection = func(collection string) MongoCollection {

	StartConnection()

	return Database.Collection(collection)
}

var getMongoClient = func() (MongoClient, error) {
	connection := getConnectionString()

	client, err := newMongoClient(connection)
	if err != nil {
		return nil, err
	}

	return client, nil
}

var getConnectionString = func() string {
	config := &settings.MongoDB

	connectionString := fmt.Sprintf("mongodb://%s:%s@%s:%s", config.User, config.Password, config.Host, config.Port)

	return connectionString
}
