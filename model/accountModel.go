package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Account struct {
	AccountID     primitive.ObjectID `bson:"_id"`
	UserID        *string            `json:"user_id"`
	AccountType   *string            `json:"account_type"`
	Balance       *float32           `json:"balance"`
	AccountStatus *string            `json:"account_status"`
	CreatedAt     int64              `json:"created_at"`
	UpdatedAt     int64              `json:"updated_at"`
}
