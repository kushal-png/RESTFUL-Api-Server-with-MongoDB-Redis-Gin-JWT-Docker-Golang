package utils

import (
	"crypto/tls"
	"project/initializers"
	models "project/model"

	"gopkg.in/gomail.v2"
)

type EmailData struct {
	Code       string
	Subject   string
}

func SendEmail(user *models.User, data *EmailData) error {
	config, _ := initializers.LoadConfig(".")

	from := config.EmailFrom
	to := user.Email
	smtpPass := config.SmtpPass
	smtpUser := config.SmtpUser
	smtpPort := config.SmtpPort
	smtpHost := config.SmtpHost

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", data.Subject)
	m.SetBody("text/plain", data.Code)

	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
