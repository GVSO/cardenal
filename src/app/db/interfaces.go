package db

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/mongo"
)

// DocumentResult is an interface for mongo.DocumentResult
type DocumentResult interface {
	Decode(v interface{}) error
}

// MongoClient is an interface for mongo.Client
type MongoClient interface {
	Connect(ctx context.Context) error
	Database(name string, opts ...*options.DatabaseOptions) *mongo.Database
}

// MongoDatabase is an interface for mongo.Database
type MongoDatabase interface {
	Collection(name string, opts ...*options.CollectionOptions) *mongo.Collection
}

// MongoCollection is an interface for mongo.Collection
type MongoCollection interface {
	InsertOne(ctx context.Context, document interface{},
		opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	FindOne(ctx context.Context, filter interface{},
		opts ...*options.FindOneOptions) *mongo.SingleResult
	FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{},
		opts ...*options.FindOneAndUpdateOptions) *mongo.SingleResult
}
