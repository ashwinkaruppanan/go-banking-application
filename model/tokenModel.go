package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Token struct {
	UserID    primitive.ObjectID `json:"user_id" bson:"user_id"`
	Token     string             `json:"token" bson:"token"`
	CreatedAt int64              `json:"created_at" bson:"created_at"`
}
