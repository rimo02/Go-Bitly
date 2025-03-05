package controller

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rimo02/url-shortener/database"
	"github.com/rimo02/url-shortener/helper"
	"github.com/rimo02/url-shortener/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"os"
	"time"
)

func ShortenTheUrl(c *gin.Context) {
	fmt.Println(c)
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

	expirationTime := time.Now().Add(time.Duration(defaultExpiry) * time.Second)
	localizedExpiration := expirationTime.Unix()

	shortUrl := helper.Shortener(6)
	resp := &models.Response{
		ID:              primitive.NewObjectID(),
		LongUrl:         req.LongUrl,
		ShortUrl:        shortUrl,
		Hits:            0,
		LastRequestTime: time.Now().Unix(),
		ExpiresAt:       localizedExpiration,
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
