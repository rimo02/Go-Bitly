package controller

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rimo02/url-shortener/database"
	"github.com/rimo02/url-shortener/helper"
	"github.com/rimo02/url-shortener/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"os"
	"time"
)

func ShortenTheUrl(c *gin.Context) {

	var req models.Request

	collection := database.GetCollection(database.Client, os.Getenv("DB_COLLECTION"))
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)

	defer cancel()

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	defaultExpiry := 24 * 60 * 60 // 1 day
	userExpiry := 0

	if req.Days > 0 {
		userExpiry += 60 * 24 * 60 * req.Days
	}
	if req.Hours > 0 {
		userExpiry += 60 * 60 * req.Hours
	}
	if req.Minutes > 0 {
		userExpiry += req.Minutes * 60
	}
	if userExpiry > 0 {
		defaultExpiry = userExpiry
	}

	expirationTime := int(time.Now().Add(time.Duration(defaultExpiry) * time.Second).Unix())

	shortUrl := helper.Shortener(6) //generate a random string of 6 length total combinations = 62^6
	resp := &models.Response{
		ID:              primitive.NewObjectID(),
		LongUrl:         req.LongUrl,
		ShortUrl:        shortUrl,
		Hits:            0,
		ExpiresAt:       expirationTime,
		LastRequestTime: int(time.Now().Unix()),
	}
	_, err := collection.InsertOne(ctx, resp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "unable to shorten the url"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"LongUrl":    req.LongUrl,
		"ShortUrl":   fmt.Sprintf("http://localhost:%s/%s", os.Getenv("PORT"), resp.ShortUrl),
		"Expires At": time.Unix(int64(resp.ExpiresAt), 0).UTC().Format(time.RFC3339),
	})

}

func GetTheUrl(c *gin.Context) {
	shortUrl := c.Param("shorturl")

	var result models.Response

	collection := database.GetCollection(database.Client, os.Getenv("DB_COLLECTION"))
	ctx, cancel := context.WithTimeout(context.TODO(), 20*time.Second)
	defer cancel()
	err := collection.FindOne(ctx, bson.M{"surl": shortUrl}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		for i := 0; i < 3; i++ {
			time.Sleep(1 * time.Second)
			err := collection.FindOne(ctx, bson.M{"surl": shortUrl}).Decode(&result)
			if err == nil {
				collection.UpdateOne(
					context.TODO(),
					bson.M{"surl": shortUrl},
					bson.M{"$inc": bson.M{"hits": 1}},
				)
				c.Redirect(http.StatusMovedPermanently, result.LongUrl)
			} else {
				continue
			}

		}
		c.JSON(http.StatusNotFound, gin.H{"message": "url not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	collection.UpdateOne(
		context.TODO(),
		bson.M{"surl": shortUrl},
		bson.M{
			"$inc": bson.M{"hits": 1},                    // Increment hits by 1
			"$set": bson.M{"lastreq": time.Now().Unix()}, // Set last request time to current Unix time
		},
	)

	c.Redirect(http.StatusMovedPermanently, result.LongUrl)
}

func Dashboard(c *gin.Context) {
	shortUrl := c.Param("shorturl")

	var result models.Response

	collection := database.GetCollection(database.Client, os.Getenv("DB_COLLECTION"))
	ctx, cancel := context.WithTimeout(context.TODO(), 20*time.Second)
	defer cancel()
	err := collection.FindOne(ctx, bson.M{"surl": shortUrl}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotFound, gin.H{"message": "url not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	remainingTime := int64(result.ExpiresAt) - time.Now().Unix()

	if remainingTime < 0 {
		c.JSON(http.StatusOK, gin.H{"message": "your url has already expired"})
		return
	}

	days := remainingTime / (24 * 3600)
	remainingTime %= (24 * 3600)
	hours := remainingTime / 3600
	remainingTime %= 3600
	minutes := remainingTime / 60

	resp := gin.H{
		"Long Url":                result.LongUrl,
		"Short Url":               result.ShortUrl,
		"Total numer of url hits": result.Hits,
		"Expires in":              fmt.Sprintf("%d Days, %d Hours, %d mins", days, hours, minutes),
		"Last requested Time":     result.LastRequestTime,
	}

	c.JSON(http.StatusOK, resp)
}

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
