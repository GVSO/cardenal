package database

import (
	"context"
	"fmt"
	"testing"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/stretchr/testify/assert"

	"github.com/gvso/cardenal/src/app/database/mocks"
	"github.com/gvso/cardenal/src/app/settings"
)

var dbConfig = settings.MongoDBConfig{
	Host:     "localhost",
	Port:     "123",
	Database: "mydb",
	User:     "user",
	Password: "pass",
}

func TestGetConnectionString(t *testing.T) {

	assert := assert.New(t)

	// Saves current settings and restore them at the end.
	oldSettings := settings.MongoDB
	defer func() { settings.MongoDB = oldSettings }()

	// Overwrites settings
	settings.MongoDB = dbConfig

	expected := fmt.Sprintf("mongodb://%s:%s@%s:%s", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port)
	connectionString := getConnectionString()

	assert.Equal(expected, connectionString)

}

func TestStartConnection(t *testing.T) {

	assert := assert.New(t)

	// Saves current settings and restores them at the end.
	oldSettings := settings.MongoDB
	defer func() { settings.MongoDB = oldSettings }()

	// Overwrites settings.
	settings.MongoDB = dbConfig

	// Saves current function and restores it at the end.
	old := getMongoClient
	defer func() { getMongoClient = old }()

	// Overwrites getMongoClient function.
	var i = 0
	var mongoClient *mocks.MongoClient

	getMongoClient = func() (MongoClient, error) {
		if i == 0 {
			mongoClient = &mocks.MongoClient{}
		}

		i++

		return mongoClient, nil
	}

	testFailedGetClientCall(assert)
	testFailedConnection(assert)
	testSuccessfulConnection(assert)
	testAlreadyConnected(assert)
}

func TestGetMongoClient(t *testing.T) {

	assert := assert.New(t)

	// Saves current function and restores it at the end.
	old := newMongoClient
	defer func() { newMongoClient = old }()

	// Overwrites newMongoClient function.
	newMongoClient = func(uri string) (*mongo.Client, error) {
		return &mongo.Client{}, nil
	}

	client, err := getMongoClient()

	assert.Equal(&mongo.Client{}, client)
	assert.Nil(err)

	// Overwrites newMongoClient function.
	newMongoClient = func(uri string) (*mongo.Client, error) {
		return nil, fmt.Errorf("could not established connection")
	}

	client, err = getMongoClient()

	assert.Nil(client)
	assert.Equal("could not established connection", err.Error())
}

// Tests startConnection when getClient fails.
func testFailedGetClientCall(assert *assert.Assertions) {
	// Saves current function and restores it at the end.
	old := getMongoClient
	defer func() { getMongoClient = old }()

	// Overwrites getMongoClient function.
	getMongoClient = func() (MongoClient, error) {
		return nil, fmt.Errorf("could not connect")
	}

	err := startConnection()

	assert.False(connected)
	assert.Equal("could not connect", err.Error())
}

// Tests startConnection when connection to database succeeds.
func testSuccessfulConnection(assert *assert.Assertions) {
	// Asserts that no connection has been performed yet.
	assert.False(connected)

	startConnection()

	// Asserts that connected is now true.
	assert.True(connected)

	client := client.(*mocks.MongoClient)

	// Asserts that Connect was called correctly.
	assert.True(client.ConnectCall.Called)
	assert.Equal(context.TODO(), client.ConnectCall.Ctx)

	// Asserts that Database was called correctly.
	assert.True(client.DatabaseCall.Called)
	assert.Equal(dbConfig.Database, client.DatabaseCall.Name)

	// Resets values for next tests.
	resetCallValues(client)
}

// Tests startConnection when connection to database does not succeed.
func testFailedConnection(assert *assert.Assertions) {
	startConnection()

	assert.False(connected)

	client := client.(*mocks.MongoClient)

	// Database should not be called.
	assert.False(client.DatabaseCall.Called)

	resetCallValues(client)
}

// Tests when connection has already been performed.
func testAlreadyConnected(assert *assert.Assertions) {
	connected = true

	startConnection()

	assert.True(connected)

	client := client.(*mocks.MongoClient)

	assert.NotNil(client)
	assert.False(client.ConnectCall.Called)
	assert.False(client.DatabaseCall.Called)
}

func resetCallValues(c *mocks.MongoClient) {
	// Resets values for next tests.
	c.ConnectCall.Called = false
	c.DatabaseCall.Called = false
}
