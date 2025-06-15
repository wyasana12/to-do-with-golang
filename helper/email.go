package helper

import (
	"fmt"
	"net/smtp"
	"to-do-list-go/config"
)

func SendVerificationEmail(to string, token string) error {
	from := config.ENV.SMTP_EMAIL
	password := config.ENV.SMTP_PASSWORD
	host := config.ENV.SMTP_HOST
	port := config.ENV.SMTP_PORT
	appurl := config.ENV.APP_URL

	auth := smtp.PlainAuth("", from, password, host)

	subject := "Subject: Email Verification \n"
	body := fmt.Sprintf("Click The Link To Verify Your Email:\n%s/api/auth/verify-email?token=%s", appurl, token)
	msg := []byte(subject + "\n" + body)

	err := smtp.SendMail(host+":"+port, auth, from, []string{to}, msg)
	return err
}
