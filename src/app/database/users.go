package database

import (
	"context"
	"fmt"

	"github.com/gvso/cardenal/src/app/database/entity"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

var collection MongoCollection

// InsertUser inserts a new user.
var InsertUser = func(user *entity.User) (interface{}, error) {

	collection = getCollection()

	// Sets the _id field value.
	(*user).ID = objectid.New()

	res, err := collection.InsertOne(nil, user)
	if err != nil {
		return nil, err
	}

	id := res.InsertedID

	return id, nil
}

// GetUserByLinkedInID returns the document containing the user with the
// provided ID.
var GetUserByLinkedInID = func(id string,
	fields ...string) (*entity.User, error) {

	collection = getCollection()

	user := entity.User{}
	filter := bson.NewDocument(bson.EC.String("linkedin_id", id))

	err := findOne(filter).Decode(&user)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &user, nil
}

var findOne = func(filter *bson.Document) DocumentResult {
	return collection.FindOne(context.Background(), filter)
}

var getCollection = func() MongoCollection {
	startConnection()

	return database.Collection("users")
}
