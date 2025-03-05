package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Request struct {
	LongUrl string `json:"lurl"`
	Hours   int    `json:"hrs"`
	Minutes int    `json:"mins"`
	Days    int    `json:"days"`
}

type Response struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	LongUrl         string             `bson:"long_url" json:"long_url"`
	ShortUrl        string             `bson:"short_url" json:"short_url"`
	Hits            int                `bson:"hits" json:"hits"`
	ExpiresAt       int64              `bson:"expires_at" json:"expires_at"`
	LastRequestTime int64              `bson:"last_request_time" json:"last_request_time"`
}

