package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Request struct {
	LongUrl string `json:"lurl"`
}

type Response struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	LongUrl         string             `bson:"lurl"`
	ShortUrl        string             `bson:"surl"`
	Hits            int                `bson:"hits"`
	ExpiresAt       int                `bson:"expat"`
	LastRequestTime int                `bson:"lastreq"`
}
