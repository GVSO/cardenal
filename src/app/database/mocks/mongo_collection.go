package mocks

import (
	"context"
	"fmt"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/insertopt"
)

// MongoCollection is the mock structure for database.MongoCollection
type MongoCollection struct {
	InsertOneCall insertOne
}

// InsertOne mocks a call to InsertOne.
func (_m *MongoCollection) InsertOne(ctx context.Context, document interface{},
	opts ...insertopt.One) (*mongo.InsertOneResult, error) {

	times := &_m.InsertOneCall.times

	// Error on second call
	if *times == 1 {
		(*times)++

		_m.InsertOneCall = insertOne{*times, true, ctx, document, opts}

		return nil, fmt.Errorf("could not insert document")
	}

	(*times)++

	_m.InsertOneCall = insertOne{*times, true, ctx, document, opts}

	return &mongo.InsertOneResult{InsertedID: "id123"}, nil
}

/*******************************************************************************
** Defines structs to check if functions were called with expected parameters **
*******************************************************************************/
type insertOne struct {
	times int

	Called   bool
	Ctx      context.Context
	Document interface{}
	Opts     []insertopt.One
}
