package auth

import (
	"context"
	"errors"
	"time"

	"github.com/golkhandani/shopWise/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type IAuthRepo interface {
	CheckIsUserRegistered(email string) error
	CheckIsUsernameTaken(username string) error
	GetAuthByEmailPass(email string, password string) (*models.Auth, error)
	CreateAuth(email string, password string) (*models.Auth, error)
	GetAuthByID(aid string) (*models.Auth, error)
}

type Repo struct {
	AuthCollection *mongo.Collection
}

func (ar Repo) CheckIsUserRegistered(email string) error {
	var existsUser models.Auth

	err := ar.AuthCollection.FindOne(context.TODO(), bson.D{{Key: "email", Value: email}}).Decode(&existsUser)
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}

	if existsUser.ID.IsZero() {
		return errors.New("user already registered")
	}
	return nil
}

func (ar Repo) CheckIsUsernameTaken(username string) error {
	var existsUser models.Auth

	err := ar.AuthCollection.FindOne(context.TODO(), bson.D{{Key: "username", Value: username}}).Decode(&existsUser)
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}

	if (err != nil && err == mongo.ErrNoDocuments) || existsUser.ID.IsZero() {
		return errors.New("username already taken")
	}
	return nil
}

func (ar Repo) GetAuthByEmailPass(email string, password string) (*models.Auth, error) {
	var existsUser models.Auth

	err := ar.AuthCollection.FindOne(context.TODO(), bson.D{{Key: "email", Value: email}}).Decode(&existsUser)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}

	if existsUser.ID.IsZero() {
		return nil, errors.New("login credentials are invalid")
	}

	matched := _CheckPasswordHash(password, existsUser.PasswordHash)
	if !matched {
		return nil, errors.New("login credentials are invalid")
	}

	return &existsUser, nil
}

func (ar Repo) CreateAuth(email string, password string) (*models.Auth, error) {
	hashedPassword, err := _HashPassword(password)
	if err != nil {
		return nil, err
	}
	dt := time.Now()

	var insertedAuth models.Auth
	auth := models.Auth{
		Email:        email,
		PasswordHash: hashedPassword,
		CreatedAt:    dt,
		UpdatedAt:    dt,
	}
	insertResult, err := ar.AuthCollection.InsertOne(context.TODO(), auth)
	if err != nil {
		return nil, err
	}

	err = ar.AuthCollection.FindOne(
		context.TODO(),
		bson.D{{Key: "_id", Value: insertResult.InsertedID}},
	).Decode(&insertedAuth)
	if err != nil {
		return nil, err
	}
	return &insertedAuth, nil
}

func (ar Repo) GetAuthByID(aid string) (*models.Auth, error) {
	var existsAuth models.Auth
	objectID, err := primitive.ObjectIDFromHex(aid)

	if err != nil {
		return nil, err
	}

	err = ar.AuthCollection.FindOne(context.TODO(), bson.D{{Key: "_id", Value: objectID}}).Decode(&existsAuth)
	if err != nil {
		return nil, err
	}

	return &existsAuth, nil
}

func _HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func _CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
