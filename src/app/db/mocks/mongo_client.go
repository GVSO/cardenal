package mocks

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoClient is the mock structure for database.MongoClient.
type MongoClient struct {
	ConnectCall  connect
	DatabaseCall database
}

// Connect mocks a call to Connect.
func (_m *MongoClient) Connect(ctx context.Context) error {

	times := &_m.ConnectCall.times

	// Error on first call.
	if *times == 0 {
		_m.ConnectCall = connect{0, true, ctx}

		(*times)++

		return errors.New("could not connect")
	}

	_m.ConnectCall = connect{*times, true, ctx}

	(*times)++

	// No errors in subsequent calls.
	return nil
}

// Database mocks a call to Database.
func (_m *MongoClient) Database(name string, opts ...*options.DatabaseOptions) *mongo.Database {
	_m.DatabaseCall = database{true, name, opts}

	return nil
}

/*******************************************************************************
** Defines structs to check if functions were called with expected parameters **
*******************************************************************************/
type connect struct {
	times int

	Called bool
	Ctx    context.Context
}
type database struct {
	Called bool
	Name   string
	opts   []*options.DatabaseOptions
}
