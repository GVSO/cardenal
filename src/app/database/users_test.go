package database

import (
	"context"
	"testing"

	"github.com/gvso/cardenal/src/app/database/mocks"
	"github.com/mongodb/mongo-go-driver/mongo/collectionopt"
	"github.com/stretchr/testify/assert"
)

var testUser = map[string]interface{}{
	"username": "username123",
}

func TestInsertUser(t *testing.T) {

	assert := assert.New(t)

	// Saves current function and restores it at the end.
	old := getCollection
	defer func() { getCollection = old }()

	// Overwrites getCollection function.
	var i = 0
	var mongoCollection *mocks.MongoCollection

	getCollection = func() MongoCollection {
		if i == 0 {
			mongoCollection = &mocks.MongoCollection{}
		}

		i++

		return mongoCollection
	}

	id, err := InsertUser(testUser)

	assert.Nil(err)
	assert.Equal("id123", id)

	collection := collection.(*mocks.MongoCollection)

	assert.True(collection.InsertOneCall.Called)
	assert.Equal(context.Background(), collection.InsertOneCall.Ctx)

	// This time, call to InsertOne returns an error
	id, err = InsertUser(testUser)
	assert.Nil(id)
	assert.Equal("could not insert document", err.Error())
}

func TestGetCollection(t *testing.T) {

	assert := assert.New(t)

	// Saves current function and restores it at the end.
	old := startConnection
	defer func() { startConnection = old }()

	// Overwrites startConnection function.
	startConnection = func() error {
		database = &mocks.MongoDatabase{}

		return nil
	}

	getCollection()

	database := database.(*mocks.MongoDatabase)

	assert.True(database.CollectionCall.Called)
	assert.Equal("users", database.CollectionCall.Name)
	// Nil options.
	assert.Equal([]collectionopt.Option([]collectionopt.Option(nil)),
		database.CollectionCall.Opts)
}
