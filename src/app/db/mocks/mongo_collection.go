package mocks

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoCollection is the mock structure for database.MongoCollection
type MongoCollection struct {
	FindOneCall          findOne
	FindOneAndUpdateCall findOneAndUpdate
	InsertOneCall        insertOne
}

// FindOne mocks a call to FindOne.
func (_m *MongoCollection) FindOne(ctx context.Context, filter interface{},
	opts ...*options.FindOneOptions) *mongo.SingleResult {

	_m.FindOneCall = findOne{true, ctx, filter, opts}

	return nil
}

// FindOneAndUpdate mocks a call to FindOneAndUpdate.
func (_m *MongoCollection) FindOneAndUpdate(ctx context.Context,
	filter interface{}, update interface{},
	opts ...*options.FindOneAndUpdateOptions) *mongo.SingleResult {

	_m.FindOneAndUpdateCall = findOneAndUpdate{true, ctx, filter, update, opts}

	return nil
}

// InsertOne mocks a call to InsertOne.
func (_m *MongoCollection) InsertOne(ctx context.Context, document interface{},
	opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {

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
	Opts    []*options.FindOneOptions
}
type findOneAndUpdate struct {
	Called bool
	Ctx    context.Context
	Filter interface{}
	Update interface{}
	Opts   []*options.FindOneAndUpdateOptions
}
type insertOne struct {
	times int

	Called   bool
	Ctx      context.Context
	Document interface{}
	Opts     []*options.InsertOneOptions
}
