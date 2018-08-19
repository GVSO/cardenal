package database

import (
	"context"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/collectionopt"
	"github.com/mongodb/mongo-go-driver/mongo/dbopt"
	"github.com/mongodb/mongo-go-driver/mongo/findopt"
	"github.com/mongodb/mongo-go-driver/mongo/insertopt"
)

// MongoClient is an interface for mongo.Client
type MongoClient interface {
	Connect(ctx context.Context) error
	Database(name string, opts ...dbopt.Option) *mongo.Database
}

// MongoDatabase is an interface for mongo.Database
type MongoDatabase interface {
	Collection(name string, opts ...collectionopt.Option) *mongo.Collection
}

// MongoCollection is an interface for mongo.Collection
type MongoCollection interface {
	InsertOne(ctx context.Context, document interface{},
		opts ...insertopt.One) (*mongo.InsertOneResult, error)
	FindOne(ctx context.Context, filter interface{},
		opts ...findopt.One) *mongo.DocumentResult
}
