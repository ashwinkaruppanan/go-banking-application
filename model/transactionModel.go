package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Transaction struct {
	TransactionID   primitive.ObjectID `bson:"_id"`
	TransactionTime int64              `json:"transaction_time"`
	SRC_Account     *string            `json:"src_account"`
	DES_Account     *string            `json:"des_account"`
	TXN_Source      *string            `json:"txn_source"`
	Operation       *string            `json:"operation"`
	Amount          *float32           `json:"amount"`
}
