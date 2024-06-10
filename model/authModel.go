package models

// ? SignUpInput struct
type SignUpInput struct {
	Name             string `json:"name" binding:"required"`
	Email            string `json:"email" binding:"required"`
	Password         string `json:"password" binding:"required"`
	PasswordConfirm  string `json:"passwordConfirm" binding:"required"`
	VerificationCode string `json:"verificationCode,omitempty"`
}

// ? SignInInput struct
type SignInInput struct {
	Email    string `json:"email" bson:"email" binding:"required"`
	Password string `json:"password" bson:"password" binding:"required"`
}

type ForgotPasswordInput struct {
	Email string `json:"email" bson:"email" binding:"required"`
}

type ResetPasswordInput struct {
	Password        string `json:"password" bson:"password" binding:"required,min=8"`
	PasswordConfirm string `json:"passwordConfirm" bson:"passwordConfirm,omitempty" binding:"required"`
}
