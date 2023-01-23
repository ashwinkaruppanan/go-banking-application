package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Transaction struct {
	TransactionID   primitive.ObjectID `bson:"_id"`
	UserID          primitive.ObjectID `json:"user_id" bson:"user_id"`
	TransactionTime int64              `json:"transaction_time" bson:"transaction_time"`
	SRC_Account     *string            `json:"src_account" bson:"src_account"`
	DES_Account     *string            `json:"des_account" bson:"des_account"`
	TXN_Source      *string            `json:"txn_source" bson:"txn_source"`
	Operation       *string            `json:"operation" bson:"operation"`
	Amount          *float32           `json:"amount" bson:"amount"`
}
