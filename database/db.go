package database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

var Client *mongo.Client

func ConnectDB() *mongo.Client {
	var uri = os.Getenv("DB_HOST")

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://localhost:%s", uri))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Could not connect to MongoDB: ", err)
	}

	fmt.Printf("Connected to MongoDB")
	Client = client
	return client
}

func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	return client.Database("url-shortener").Collection(collectionName)
}
