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
func DbInstance() *mongo.Client {

	client, err := mongo.NewClient(options.Client().ApplyURI(connectionString))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	err = client.Connect(ctx)

	if err != nil {

		log.Fatal(err)
	}
	fmt.Println("Connected to DB ")
	return client
}

var Client *mongo.Client = DbInstance()

func OpenCollection(client *mongo.Client, collection_Name string) *mongo.Collection {
	collection := client.Database(dbName).Collection(collection_Name)
	return collection

}
