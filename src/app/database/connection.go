package database

import (
	"context"
	"fmt"
	"log"

	"github.com/mongodb/mongo-go-driver/mongo"

	"github.com/gvso/cardenal/src/app/settings"
)

var connected = false

var client *mongo.Client

var database *mongo.Database

func startConnection() {

	if !connected {
		var err error

		connection := getConnectionString()

		fmt.Println(connection)

		client, err = mongo.NewClient(connection)
		if err != nil {
			log.Fatal(err)
		}

		err = client.Connect(context.TODO())
		if err != nil {
			log.Fatal(err)
		}

		database = client.Database("cardenal")

		connected = true
	}
}

func getConnectionString() string {
	config := settings.MongoDB

	connectionString := "mongodb://"
	connectionString += config.User + ":" + config.Password + "@" + config.Host + ":" + config.Port

	return connectionString
}
