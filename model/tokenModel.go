package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Token struct {
	UserID    primitive.ObjectID `json:"user_id"`
	Token     string             `json:"token"`
	CreatedAt int64              `json:"created_at"`
}
