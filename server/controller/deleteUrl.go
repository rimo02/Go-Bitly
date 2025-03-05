package controller

import (
	"context"
	"github.com/rimo02/url-shortener/database"
	"github.com/rimo02/url-shortener/models"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"time"
	"os"
)


func DeleteExpiredUrls() {
	collection := database.GetCollection(database.Client, os.Getenv("DB_COLLECTION"))
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal("Error in fetching documents ", err.Error())
		return
	}
	defer cursor.Close(ctx)

	currTime := time.Now().Unix()
	for cursor.Next(ctx) {
		var item models.Response
		if err := cursor.Decode(&item); err != nil {
			log.Fatal("Error in decoding document: ", err.Error())
			continue
		}
		if currTime > int64(item.ExpiresAt) {
			_, err := collection.DeleteOne(ctx, bson.M{"_id": item.ID})
			if err != nil {
				log.Printf("Error in deleting the url")
			} else {
				log.Printf("Succesfully deleted the url")
			}
		}
	}
	if err := cursor.Err(); err != nil {
		log.Fatal("Error in iterating through cursor: ", err.Error())
		return
	}

}
