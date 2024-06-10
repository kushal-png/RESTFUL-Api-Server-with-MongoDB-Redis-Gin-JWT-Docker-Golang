package services

import models "project/model"

type AuthService interface {
	SignUpUser(*models.SignUpInput) (*models.User, error)
	VerifyUser(string)(error)
	ForgotPassword(string, string)(error)
	ResetPassword(string, string)(error)
}
