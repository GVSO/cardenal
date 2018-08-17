package database

import (
	"context"
	"fmt"

	"github.com/mongodb/mongo-go-driver/mongo"

	"github.com/gvso/cardenal/src/app/settings"
)

var connected = false

var client MongoClient

var database MongoDatabase

/**
 * Helper functions from external packages.
 *
 * For testing purposes, variables are declared which are defined as pointers to
 * function. The advantage of doing this is that these variables can later be
 * overwritten in the testing files.
 */
var newMongoClient = mongo.NewClient

var startConnection = func() error {

	if !connected {
		var err error

		client, err = getMongoClient()
		if err != nil {
			return err
		}

		err = client.Connect(context.TODO())
		if err != nil {
			return err
		}

		database = client.Database(settings.MongoDB.Database)

		connected = true
	}

	return nil
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
