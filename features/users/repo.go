package users

import (
	"context"
	"errors"

	"github.com/golkhandani/shopWise/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type IUserRepo interface {
	GetUserByAuthID(id primitive.ObjectID) (*models.User, error)
	CreateUser(authId primitive.ObjectID, username string, fullName string) (*models.User, error)
}

type UserRepo struct {
	UserCollection *mongo.Collection
}

var ErrUserNotFound = errors.New("user not found")
var ErrorUserAlreadyHaveProfile = errors.New("user already have profile")

func (ur UserRepo) GetUserByAuthID(id primitive.ObjectID) (*models.User, error) {

	var existsUser models.User

	err := ur.UserCollection.FindOne(context.TODO(), bson.D{{Key: "authId", Value: id}}).Decode(&existsUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	if existsUser.ID.IsZero() {
		return nil, ErrUserNotFound
	}

	return &existsUser, nil
}

func (ur UserRepo) CreateUser(authId primitive.ObjectID, username string, fullName string) (*models.User, error) {
	existsUser, err := ur.GetUserByAuthID(authId)

	// if user found it means user cannot create a profile again
	// they can use update instead
	// if err is user not found it means user can continue the
	// flow and create a new profile
	if err != nil && err != ErrUserNotFound {
		return nil, err
	}

	if existsUser != nil {
		return nil, ErrorUserAlreadyHaveProfile
	}
	newUser := models.NewUser(authId, username, fullName)

	inserted, err := ur.UserCollection.InsertOne(context.TODO(), newUser)
	if err != nil {
		return nil, err
	}

	var insertedUser models.User

	err = ur.UserCollection.FindOne(context.TODO(), bson.D{{Key: "_id", Value: inserted.InsertedID}}).Decode(&insertedUser)
	if err != nil {
		return nil, err
	}

	return &insertedUser, nil
}
