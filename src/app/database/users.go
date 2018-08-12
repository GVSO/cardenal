package database

import (
	"context"
	"fmt"
	"log"

	"github.com/mongodb/mongo-go-driver/mongo"
)

var collection *mongo.Collection

// InsertUser inserts a new user.
func InsertUser(user map[string]interface{}) {

	collection = getCollection()

	res, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		log.Fatal(err)
	}

	id := res.InsertedID

	fmt.Println(id)
}

func getCollection() *mongo.Collection {
	startConnection()

	return database.Collection("users")
}
