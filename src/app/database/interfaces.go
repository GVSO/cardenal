package database

import (
	"context"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/collectionopt"
	"github.com/mongodb/mongo-go-driver/mongo/dbopt"
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
