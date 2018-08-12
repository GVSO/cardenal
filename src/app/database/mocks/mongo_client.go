package mocks

import (
	"context"
	"fmt"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/dbopt"
)

// MongoClient is the mock structure for database.MongoClient
type MongoClient struct {
	ConnectCall  connect
	DatabaseCall database
}

var i = -1

// Connect mocks a call to Connect.
func (_m *MongoClient) Connect(ctx context.Context) error {

	_m.ConnectCall = connect{true, ctx}

	i++
	// Error on first call.
	if i == 0 {
		return fmt.Errorf("could not connect")
	}

	// No errors in subsequent calls
	return nil
}

// Database mocks a call to Database.
func (_m *MongoClient) Database(name string, opts ...dbopt.Option) *mongo.Database {
	_m.DatabaseCall = database{true, name, opts}

	return nil
}

/*******************************************************************************
** Defines structs to check if functions were called with expected parameters **
*******************************************************************************/
type connect struct {
	Called bool
	Ctx    context.Context
}
type database struct {
	Called bool
	Name   string
	opts   []dbopt.Option
}
