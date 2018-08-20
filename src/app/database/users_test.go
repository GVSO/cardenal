package database

import (
	"testing"

	"github.com/gvso/cardenal/src/app/database/entity"

	"github.com/gvso/cardenal/src/app/database/mocks"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo/collectionopt"
	"github.com/stretchr/testify/assert"
)

var testUser = entity.User{
	LinkedInID: "linkedin_id123",
	FirstName:  "John",
	LastName:   "Smith",
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

	id, err := InsertUser(&testUser)

	assert.Nil(err)
	assert.Equal("id123", id)

	collection := collection.(*mocks.MongoCollection)

	assert.True(collection.InsertOneCall.Called)
	assert.Equal(&testUser, collection.InsertOneCall.Document)

	// This time, call to InsertOne returns an error
	id, err = InsertUser(&testUser)
	assert.Nil(id)
	assert.Equal("could not insert document", err.Error())
}

func TestGetUserByLinkedInID(t *testing.T) {

	assert := assert.New(t)

	// Saves current function and restores it at the end.
	old := getCollection
	defer func() { getCollection = old }()

	// Overwrites getCollection function.
	getCollection = func() MongoCollection {
		return &mocks.MongoCollection{}
	}

	// Saves current function and restores it at the end.
	oldFindOne := findOne
	defer func() { findOne = oldFindOne }()

	var filter *bson.Document
	dr := &mocks.DocumentResult{}

	// Overwrites findOne function.
	findOne = func(filterArg *bson.Document) DocumentResult {
		filter = filterArg
		return dr
	}

	user, err := GetUserByLinkedInID("linkedin_id123")

	// Asserts the value of the filter.
	assert.Equal("linkedin_id", filter.ElementAt(0).Key())
	assert.Equal("linkedin_id123", filter.ElementAt(0).Value().StringValue())

	assert.True(dr.DecodeCall.Called)
	assert.Equal(entity.User{}, dr.DecodeCall.Value)

	expected := entity.User{
		LinkedInID: "linkedin_id123",
		FirstName:  "John",
		LastName:   "Smith",
	}

	assert.Nil(err)
	assert.Equal(expected, *user)

	// mocks.DocumentResult.Decode() returns an error the second time
	user, err = GetUserByLinkedInID("linkedin_id123")

	assert.Nil(user)
	assert.Equal("could not decode document", err.Error())
}

func TestFindOne(t *testing.T) {

	assert := assert.New(t)

	// Overwrites collection value.
	collection = &mocks.MongoCollection{}

	filter := bson.NewDocument(bson.EC.String("linkedin_id", "linkedin_id123"))

	findOne(filter)

	collection := collection.(*mocks.MongoCollection)

	assert.True(collection.FindOneCall.Callled)
	assert.Equal(filter, collection.FindOneCall.Filter)
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
