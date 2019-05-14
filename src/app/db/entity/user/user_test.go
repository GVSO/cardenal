package user

import (
	"testing"

	"github.com/gvso/cardenal/src/app/db"
	"github.com/gvso/cardenal/src/app/db/mocks"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

var testUser = User{
	LinkedInID: "linkedin_id123",
	FirstName:  "John",
	LastName:   "Smith",
}

func TestInsertUser(t *testing.T) {

	assert := assert.New(t)

	// Saves current collection and restores it at the end.
	old := collection
	defer func() { collection = old }()

	collection = &mocks.MongoCollection{}

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

	// Saves current collection and restores it at the end.
	old := collection
	defer func() { collection = old }()

	collection = &mocks.MongoCollection{}

	// Saves current function and restores it at the end.
	oldFindOne := findOne
	defer func() { findOne = oldFindOne }()

	dr := &mocks.SingleResult{}

	// Overwrites findOne function.
	findOne = func(filter interface{}) db.SingleResult {
		// Asserts the value of the filter.
		expected := bson.D{{"linkedin_id", "linkedin_id123"}}
		f := filter.(bson.D)

		assert.Equal(expected, f)

		return dr
	}

	user, err := GetUserByLinkedInID("linkedin_id123")

	assert.True(dr.DecodeCall.Called)
	assert.Equal(User{}, dr.DecodeCall.Value)

	expected := User{
		LinkedInID: "linkedin_id123",
		FirstName:  "John",
		LastName:   "Smith",
	}

	assert.Nil(err)
	assert.Equal(expected, *user)

	// mocks.SingleResult.Decode() returns an error the second time
	user, err = GetUserByLinkedInID("linkedin_id123")

	assert.Nil(user)
	assert.Equal("could not decode document", err.Error())
}

func TestUpdateUserByLinkedInID(t *testing.T) {

	assert := assert.New(t)

	// Saves current collection and restores it at the end.
	old := collection
	defer func() { collection = old }()

	collection = &mocks.MongoCollection{}

	// Saves current function and restores it at the end.
	oldFindOneAndUpdate := findOneAndUpdate
	defer func() { findOneAndUpdate = oldFindOneAndUpdate }()

	update := map[string]interface{}{
		"$set": map[string]string{
			"access_token": "access_token123",
		},
	}

	dr := &mocks.SingleResult{}

	// Overwrites findOne function.
	findOneAndUpdate = func(filterArg *bson.D, updateArg interface{}) db.SingleResult {
		// Asserts the value of the filter.
		expected := &bson.D{{"linkedin_id", "linkedin_id123"}}
		assert.Equal(expected, filterArg)

		// Asserts the value of update.
		assert.Equal(update, updateArg)

		return dr
	}

	user, err := UpdateUserByLinkedInID("linkedin_id123", update)

	// Asserts the values passed to Decode.
	assert.True(dr.DecodeCall.Called)
	assert.Equal(User{}, dr.DecodeCall.Value)

	// Asserts the values returned by UpdateUserByLinkedInUpdate.
	expected := User{
		LinkedInID: "linkedin_id123",
		FirstName:  "John",
		LastName:   "Smith",
	}

	assert.Nil(err)
	assert.Equal(expected, *user)

	// mocks.SingleResult.Decode() returns an error the second time
	user, err = UpdateUserByLinkedInID("linkedin_id123", update)

	assert.Nil(user)
	assert.Equal("could not decode document", err.Error())
}

func TestIsTokenValid(t *testing.T) {
	assert := assert.New(t)

	// Saves current collection and restores it at the end.
	old := collection
	defer func() { collection = old }()

	collection = &mocks.MongoCollection{}

	// Saves current function and restores it at the end.
	oldFindOne := findOne
	defer func() { findOne = oldFindOne }()

	dr := &mocks.SingleResult{}

	// Overwrites findOne function.
	findOne = func(filter interface{}) db.SingleResult {
		// Asserts the value of the filter.
		f := filter.(map[string]string)
		assert.Equal("linkedin_id123", f["linkedin_id"])
		assert.Equal("token123", f["access_token"])

		return dr
	}

	isTokenValid := IsTokenValid("linkedin_id123", "token123")

	assert.True(dr.DecodeCall.Called)
	assert.Equal(User{}, dr.DecodeCall.Value)

	assert.True(isTokenValid)

	// mocks.SingleResult.Decode() returns an error the fourth time
	isTokenValid = IsTokenValid("linkedin_id123", "token123")

	assert.False(isTokenValid)
}

func TestFindOne(t *testing.T) {

	assert := assert.New(t)

	// Overwrites collection value.
	collection = &mocks.MongoCollection{}

	filter := bson.D{{"linkedin_id", "linkedin_id123"}}

	findOne(filter)

	collection := collection.(*mocks.MongoCollection)

	assert.True(collection.FindOneCall.Callled)
	assert.Equal(filter, collection.FindOneCall.Filter)
}

func TestFindOneAndUpdate(t *testing.T) {

	assert := assert.New(t)

	// Overwrites collection value.
	collection = &mocks.MongoCollection{}

	filter := &bson.D{{"linkedin_id", "linkedin_id123"}}
	update := &bson.D{{
		"$set", bson.E{
			"access_token", "access_token123",
		},
	}}

	findOneAndUpdate(filter, update)

	collection := collection.(*mocks.MongoCollection)

	assert.True(collection.FindOneAndUpdateCall.Called)
	assert.Equal(filter, collection.FindOneAndUpdateCall.Filter)
	assert.Equal(update, collection.FindOneAndUpdateCall.Update)
}
