package services

import (
	"context"
	"errors"
	models "project/model"
	"project/utils"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AuthServiceImpl struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewAuthServiceImpl(cl *mongo.Collection, ct context.Context) AuthServiceImpl {
	return AuthServiceImpl{
		collection: cl,
		ctx:        ct,
	}
}

func (a *AuthServiceImpl) SignUpUser(signUpInput *models.SignUpInput) (*models.User, error) {
	//Intializing the user from signUpInput
	user := &models.User{
		ID:               primitive.NewObjectID(),
		Name:             signUpInput.Name,
		Email:            strings.ToLower(signUpInput.Email),
		Password:         utils.HashPassword(signUpInput.Password),
		Role:             "user",
		Verified:         false,
		VerificationCode: signUpInput.VerificationCode,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	// Insert the user into the collection
	_, err := a.collection.InsertOne(a.ctx, user)
	if err != nil {
		return nil, errors.New("failed to create new user")
	}

	// Create a unique index for the email field
	opt := options.Index()
	opt.SetUnique(true)
	indexModel := mongo.IndexModel{
		Keys:    bson.M{"email": 1}, // 1 for ascending order, -1 for descending order
		Options: opt,
	}

	if _, err := a.collection.Indexes().CreateOne(a.ctx, indexModel); err != nil {
		return nil, errors.New("email already exists")
	}

	return user, nil
}

func (a *AuthServiceImpl) VerifyUser(code string) error {
	query := bson.D{{Key: "verificationCode", Value: code}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "verified", Value: true}}}}
	res, err := a.collection.UpdateOne(a.ctx, query, update)

	if err != nil {
		return errors.New("failed to update")
	}
	if res.MatchedCount == 0 {
		return errors.New("invalid Verification code")
	}
	return nil
}

func (a *AuthServiceImpl) ForgotPassword(email string, code string) error {
	query := bson.D{{Key: "email", Value: email}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "passwordResetCode", Value: code}}}}
	res, err := a.collection.UpdateOne(a.ctx, query, update)

	if err != nil {
		return errors.New("failed to update")
	}
	if res.MatchedCount == 0 {
		return errors.New("invalid email")
	}
	return nil
}

func (a *AuthServiceImpl) ResetPassword(hashedPw string, code string) error {
	query := bson.D{{Key: "passwordResetCode", Value: code}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "password", Value: hashedPw}}}, {Key: "$unset", Value: bson.D{{Key: "passwordResetCode", Value: ""}}}}

	res, err := a.collection.UpdateOne(a.ctx, query, update)
	if err != nil {
		return errors.New("failed to update")
	}
	if res.MatchedCount == 0 {
		return errors.New("invalid reset code")
	}
	return nil
}
