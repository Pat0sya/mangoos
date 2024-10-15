package configs

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ConnectDB establishes a connection to the MongoDB database and returns a client.
func ConnectDB() *mongo.Client {
	// Set up the MongoDB client options using the URI from the environment
	clientOptions := options.Client().ApplyURI(EnvMongoURI())

	// Use the context with timeout for connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to the MongoDB server
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
	return client
}

var DB *mongo.Client = ConnectDB()

// GetCollection returns a reference to a MongoDB collection.
func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("golangAPI").Collection(collectionName)
	return collection
}
