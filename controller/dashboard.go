package controller

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rimo02/url-shortener/database"
	"github.com/rimo02/url-shortener/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"os"
	"time"
)

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
