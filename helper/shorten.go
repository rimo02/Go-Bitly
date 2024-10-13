package helper

import (
	"context"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/rimo02/url-shortener/database"
	"github.com/rimo02/url-shortener/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const letter = "AqSw0eQ1rt2yW3Sui4DoEp5Fas6GR7fHgJ8hKTjkLYlU9zIZxOXvbPVnNmCBNM"

func Shortener(n int) string {
	collection := database.GetCollection(database.Client, os.Getenv("DB_COLLECTION"))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	b := make([]byte, n)
	for {
		for i := range b {
			b[i] = letter[rand.Intn(len(letter))]
		}
		var result models.Response
		err := collection.FindOne(ctx, bson.M{"surl": string(b)}).Decode(&result)
		if err == mongo.ErrNoDocuments {
			return string(b)
		} else if err != nil {
			log.Printf("A database error occured")
		} else {
			continue
		}
	}
}
