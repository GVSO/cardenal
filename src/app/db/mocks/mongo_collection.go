package mocks

import (
	"context"
	"errors"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/findopt"
	"github.com/mongodb/mongo-go-driver/mongo/insertopt"
)

// MongoCollection is the mock structure for database.MongoCollection
type MongoCollection struct {
	FindOneCall          findOne
	FindOneAndUpdateCall findOneAndUpdate
	InsertOneCall        insertOne
}

// FindOne mocks a call to FindOne.
func (_m *MongoCollection) FindOne(ctx context.Context, filter interface{},
	opts ...findopt.One) *mongo.DocumentResult {

	_m.FindOneCall = findOne{true, ctx, filter, opts}

	return nil
}

// FindOneAndUpdate mocks a call to FindOneAndUpdate.
func (_m *MongoCollection) FindOneAndUpdate(ctx context.Context,
	filter interface{}, update interface{},
	opts ...findopt.UpdateOne) *mongo.DocumentResult {

	_m.FindOneAndUpdateCall = findOneAndUpdate{true, ctx, filter, update, opts}

	return nil
}

// InsertOne mocks a call to InsertOne.
func (_m *MongoCollection) InsertOne(ctx context.Context, document interface{},
	opts ...insertopt.One) (*mongo.InsertOneResult, error) {

	times := &_m.InsertOneCall.times

	// Error on second call
	if *times == 1 {
		(*times)++

		_m.InsertOneCall = insertOne{*times, true, ctx, document, opts}

		return nil, errors.New("could not insert document")
	}

	(*times)++

	_m.InsertOneCall = insertOne{*times, true, ctx, document, opts}

	return &mongo.InsertOneResult{InsertedID: "id123"}, nil
}

/*******************************************************************************
** Defines structs to check if functions were called with expected parameters **
*******************************************************************************/
type findOne struct {
	Callled bool
	Ctx     context.Context
	Filter  interface{}
	Opts    []findopt.One
}
type findOneAndUpdate struct {
	Called bool
	Ctx    context.Context
	Filter interface{}
	Update interface{}
	Opts   []findopt.UpdateOne
}
type insertOne struct {
	times int

	Called   bool
	Ctx      context.Context
	Document interface{}
	Opts     []insertopt.One
}
