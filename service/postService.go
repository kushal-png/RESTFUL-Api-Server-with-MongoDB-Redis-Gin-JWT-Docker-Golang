package services

import models "project/model"

type PostServices interface {
	GetPost(string) (*models.Post, error)
	GetPosts(int, int) ([]*models.Post, error)
	DeletePost(string) error
	CreatePost(*models.CreatePost) (*models.Post, error)
	UpdatePost(*models.UpdatePost, string) (*models.Post, error)
}
