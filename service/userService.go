package services

import models "project/model"

type UserService interface {
	GetUserByMail(string) (*models.User, error)
	GetUserById(string) (*models.User, error)
}
