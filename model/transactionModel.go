package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Transaction struct {
	TransactionID   primitive.ObjectID `bson:"_id"`
	UserID          primitive.ObjectID `json:"user_id" bson:"user_id"`
	TransactionTime int64              `json:"transaction_time" bson:"transaction_time"`
	SRC_Account     primitive.ObjectID `json:"src_account" bson:"src_account"`
	DES_Account     primitive.ObjectID `json:"des_account" bson:"des_account"`
	Operation       string             `json:"operation" bson:"operation"`
	Amount          float32            `json:"amount" bson:"amount"`
}

type Transfer struct {
	DES_Account string  `json:"des_account"`
	Amount      float32 `json:"amount"`
}
