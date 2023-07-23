package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dbName = "uGrow"
const colName = "equity"
const connectionString = "mongodb://127.0.0.1:27017/" + dbName

// MOST IMPORTANT
var collection *mongo.Collection

// connect with  mongoDB
func DbInstance() (*mongo.Client, error) {

	client, err := mongo.NewClient(options.Client().ApplyURI(connectionString))
	if err != nil {
		// log.Fatal(err)
		return nil, err

	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	err = client.Connect(ctx)

	if err != nil {

		log.Fatal(err)
	}
	fmt.Println("Connected to DB ")
	return client, nil
}

var Client *mongo.Client

func init() {
	// Initialize the MongoDB client when the package is imported.
	var err error
	Client, err = DbInstance()

	if err != nil {
		log.Fatal("Failed to initialize the MongoDB client:", err)
	}
	// Client = &client
}

func OpenCollection(client *mongo.Client) *mongo.Collection {
	collection := client.Database(dbName).Collection("user")
	return collection

}
func OpenWatchListCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database(dbName).Collection(collectionName)
	return collection

}
