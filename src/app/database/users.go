package database

import (
	"context"
)

var collection MongoCollection

// InsertUser inserts a new user.
var InsertUser = func(user map[string]interface{}) (interface{}, error) {

	collection = getCollection()

	res, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		return nil, err
	}

	id := res.InsertedID

	return id, nil
}

var getCollection = func() MongoCollection {
	startConnection()

	return database.Collection("users")
}
