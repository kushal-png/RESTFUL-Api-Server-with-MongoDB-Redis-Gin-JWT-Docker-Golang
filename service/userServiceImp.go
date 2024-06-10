package services

import (
	"context"
	models "project/model"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserServiceImpl struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewUserServiceImpl(cl *mongo.Collection, ct context.Context) UserServiceImpl {
	return UserServiceImpl{cl, ct}
}

func (u *UserServiceImpl) GetUserByMail(email string) (*models.User, error) {
	var res *models.User
	query := bson.M{"email": strings.ToLower(email)}
	err := u.collection.FindOne(u.ctx, query).Decode(&res)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &models.User{}, err
		}
		return nil, err
	}

	return res, nil
}

func (u *UserServiceImpl) GetUserById(Id string) (*models.User, error) {
	oid, _ := primitive.ObjectIDFromHex(Id)

	var user *models.User

	query := bson.M{"_id": oid}
	err := u.collection.FindOne(u.ctx, query).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &models.User{}, err
		}
		return nil, err
	}

	return user, nil
}
