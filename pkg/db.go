package pkg

import (
	"context"
	"fmt"
	"log"

	"github.com/SwishHQ/spread/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ConnectDB initializes the database connection
func MongoConnection() (*mongo.Database, error) {
	// Set client options
	clientOptions := options.Client().ApplyURI(config.MongoUrl)
	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal("Error", err)
	}
	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
	return client.Database(config.MongoDatabase), nil
}
