package mocks

import (
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/collectionopt"
)

// MongoDatabase is the mock structure for database.MongoDatabase
type MongoDatabase struct {
	CollectionCall collection
}

// Collection mocks a call to Collection.
func (_m *MongoDatabase) Collection(name string, opts ...collectionopt.Option) *mongo.Collection {
	_m.CollectionCall = collection{true, name, opts}

	return nil
}

/*******************************************************************************
** Defines structs to check if functions were called with expected parameters **
*******************************************************************************/
type collection struct {
	Called bool
	Name   string
	Opts   []collectionopt.Option
}
