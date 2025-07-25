package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Username     string             `json:"username" bson:"username" validate:"required,min=3,max=20"`
	Password     string             `json:"password" bson:"password" validate:"required,min=6,max=100"`
	RefreshToken string             `json:"token" bson:"token"`
	UserType     string             `json:"user_type" bson:"user_type" validate:"required,oneof=ADMIN USER"` // "ADMIN" or "USER"
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
	UserID       string             `json:"user_id" bson:"user_id" validate:"required,min=1,max=100"`
}

// for user login
type Crediential struct {
	UserID   string `json:"user_id" validate:"required,min=1,max=100"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
