package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Environement struct {
	MongodbURI   string
	DatabaseName string
	JWTSecret    string
	Port         int64
}

type Auth struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Email        string             `json:"email" bson:"email"`
	PasswordHash string             `json:"-" bson:"passwordHash"`
	CreatedAt    time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt    time.Time          `json:"updatedAt" bson:"updatedAt"`
}

type User struct {
	// bson:"_id,omitempty" to remove it when we creating an item
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	// json:"-" as we don't want to expose auth details
	AuthID    primitive.ObjectID `json:"-" bson:"authId"`
	Username  string             `json:"username" bson:"username"`
	FullName  string             `json:"fullName" bson:"fullName"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}

func NewUser(authID primitive.ObjectID, username string, fullName string) User {
	dt := time.Now()
	return User{
		AuthID:    authID,
		Username:  username,
		FullName:  fullName,
		CreatedAt: dt,
		UpdatedAt: dt,
	}
}
