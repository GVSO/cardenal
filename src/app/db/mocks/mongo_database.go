package mocks

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDatabase is the mock structure for database.MongoDatabase
type MongoDatabase struct {
	CollectionCall collection
}

// Collection mocks a call to Collection.
func (_m *MongoDatabase) Collection(name string, opts ...*options.CollectionOptions) *mongo.Collection {
	_m.CollectionCall = collection{true, name, opts}

	return nil
}

/*******************************************************************************
** Defines structs to check if functions were called with expected parameters **
*******************************************************************************/
type collection struct {
	Called bool
	Name   string
	Opts   []*options.CollectionOptions
}
