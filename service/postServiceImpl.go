package services

import (
	"context"
	"errors"
	models "project/model"
	"project/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PostServiceImpl struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewPostServiceImpl(col *mongo.Collection, c context.Context) PostServiceImpl {
	return PostServiceImpl{
		collection: col,
		ctx:        c,
	}
}

func (p *PostServiceImpl) GetPost(id string) (*models.Post, error) {
	var response *models.Post
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid ID format")
	}

	filter := bson.D{{Key: "_id", Value: objectID}}
	err = p.collection.FindOne(p.ctx, filter).Decode(&response)
	if err != nil {
		return nil, errors.New("cannot find post with given id")
	}
	return response, nil
}

func (p *PostServiceImpl) GetPosts(page int, limit int) ([]*models.Post, error) {
	if page == 0 {
		page = 1
	}

	if limit == 0 {
		limit = 10
	}

	skip := (page - 1) * limit
	opt := options.FindOptions{}
	opt.SetLimit(int64(limit))
	opt.SetSkip(int64(skip))

	query := bson.M{}
	cursor, err := p.collection.Find(p.ctx, query, &opt)
	if err != nil {
		return nil, errors.New("cannot find posts")
	}

	defer cursor.Close(p.ctx)

	var response []*models.Post
	for cursor.Next(p.ctx) {
		post := &models.Post{}
		err := cursor.Decode(post)
		if err != nil {
			return nil, errors.New("error in decoding the cursor")
		}
		response = append(response, post)
	}

	if len(response) == 0 {
		return []*models.Post{}, nil
	}
	return response, nil
}

func (p *PostServiceImpl) DeletePost(id string) error {
	object_id, _ := primitive.ObjectIDFromHex(id)
	query := bson.M{"_id": object_id}
	res, err := p.collection.DeleteOne(p.ctx, query)
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("not found")
	}
	return nil
}

func (p *PostServiceImpl) CreatePost(post *models.CreatePost) (*models.Post, error) {
	newPost := &models.Post{
		ID:        primitive.NewObjectID(),
		Title:     post.Title,
		Image:     post.Image,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Content:   post.Content,
		User:      post.User,
	}

	_, err := p.collection.InsertOne(p.ctx, newPost)
	if err != nil {
		return nil, errors.New("failed to create new post")
	}

	return newPost, nil
}

func (p *PostServiceImpl) UpdatePost(post *models.UpdatePost, id string) (*models.Post, error) {
	doc, err := utils.ToDoc(post)
	if err != nil {
		return nil, errors.New("failed to document the data")
	}
	doc = append(doc, bson.E{Key: "updatedAt", Value: time.Now()})

	object_id, _ := primitive.ObjectIDFromHex(id)
	query := bson.D{{Key: "_id", Value: object_id}}
	update := bson.D{{Key: "$set", Value: doc}}

	var updatedPost *models.Post
	err = p.collection.FindOneAndUpdate(p.ctx, query, update, options.FindOneAndUpdate().SetReturnDocument(1)).Decode(&updatedPost)
	if err != nil {
		return nil, errors.New("failed to update")
	}

	return updatedPost, nil
}
