package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Account struct {
	AccountID     primitive.ObjectID `bson:"_id"`
	UserID        primitive.ObjectID `json:"user_id" bson:"user_id"`
	AccountType   string             `json:"account_type" bson:"account_type"`
	Balance       float32            `json:"balance" bson:"balance"`
	AccountStatus string             `json:"account_status" bson:"account_status"`
	CreatedAt     int64              `json:"created_at" bson:"created_at"`
	UpdatedAt     int64              `json:"updated_at" bson:"updated_at"`
}

type ShowAccount struct {
	AccountType   string
	Balance       float32
	AccountStatus string
}

type ActivateAccount struct {
	AccountID string `json:"account_id"`
	Operation string `json:"operation"`
}
