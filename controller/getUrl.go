package controller

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"github.com/rimo02/url-shortener/database"
	"github.com/rimo02/url-shortener/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetTheUrl(c *gin.Context) {
	shortUrl := c.Param("shorturl")

	if cacheResult, found := database.Cache.Get(shortUrl); found {
		result := cacheResult.(string)
		go incrementAndUpdateDB(shortUrl)
		c.Redirect(http.StatusMovedPermanently, result)
		return
	}

	var result models.Response

	collection := database.GetCollection(database.Client, os.Getenv("DB_COLLECTION"))
	ctx, cancel := context.WithTimeout(context.TODO(), 20*time.Second)
	defer cancel()
	err := collection.FindOne(ctx, bson.M{"surl": shortUrl}).Decode(&result)
	if err == mongo.ErrNoDocuments { // add retry mechanism
		for i := 0; i < 3; i++ {
			time.Sleep(1 * time.Second)
			err := collection.FindOne(ctx, bson.M{"surl": shortUrl}).Decode(&result)
			if err == nil {
				go incrementAndUpdateDB(shortUrl)
				c.Redirect(http.StatusMovedPermanently, result.LongUrl)
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"message": "url not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	incrementAndUpdateDB(shortUrl)

	if result.Hits > database.MaxHits {
		database.Cache.Set(shortUrl, result.LongUrl, cache.DefaultExpiration)
	}

	c.Redirect(http.StatusMovedPermanently, result.LongUrl)
}

func incrementAndUpdateDB(shortUrl string) {
	collection := database.GetCollection(database.Client, os.Getenv("DB_COLLECTION"))
	collection.UpdateOne(
		context.TODO(),
		bson.M{"surl": shortUrl},
		bson.M{
			"$inc": bson.M{"hits": 1},                    // Increment hits by 1
			"$set": bson.M{"lastreq": time.Now().Unix()}, // Set last request time to current Unix time
		},
	)
}
