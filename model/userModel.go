package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	UserID     primitive.ObjectID `bson:"_id"`
	Email      string             `json:"email" validate:"email,required"`
	Password   string             `json:"password" validate:"required"`
	FullName   string             `json:"full_name" validate:"required"`
	UserType   string             `json:"user_type" validate:"required,eq=ADMIN|eq=USER"`
	UserStatus int                `json:"user_status"`
	CreatedAt  int64              `json:"created_at"`
	UpdatedAt  int64              `json:"updated_at"`
}