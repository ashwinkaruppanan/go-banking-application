package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	UserID     primitive.ObjectID `bson:"_id"`
	Email      string             `json:"email" validate:"email,required" bson:"email"`
	Password   string             `json:"password" validate:"required" bson:"password"`
	FullName   string             `json:"full_name" validate:"required" bson:"full_name"`
	UserType   string             `json:"user_type" validate:"required,eq=ADMIN|eq=USER" bson:"user_type"`
	UserStatus int                `json:"user_status" bson:"user_status"`
	CreatedAt  int64              `json:"created_at" bson:"created_at"`
	UpdatedAt  int64              `json:"updated_at" bson:"updated_at"`
}

type ShowUser struct {
	Email      string
	FullName   string
	UserStatus int
}

type UpdateUser struct {
	FullName        string `json:"full_name,omitempty"`
	Email           string `json:"email,omitempty"`
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password,omitempty"`
}
