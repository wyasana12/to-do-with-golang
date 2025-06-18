package helper

import (
	"fmt"
	"net/smtp"
	"to-do-list-go/config"
	"to-do-list-go/models"
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

func SendTodoReminderEmail(to string, todo *models.Todo, reminderType string) error {
	from := config.ENV.SMTP_EMAIL
	password := config.ENV.SMTP_PASSWORD
	host := config.ENV.SMTP_HOST
	port := config.ENV.SMTP_PORT

	auth := smtp.PlainAuth("", from, password, host)

	subject := fmt.Sprintf("Subject: Task Reminder %s\n", todo.Title)

	var body string

	switch reminderType {
	case "D-1":
		body = fmt.Sprintf("Hello,\n\nYour task '%s' is due in less than 24 hours. Don't forget to complete it!\n\nDeadline: %s\n", todo.Title, todo.EndDate.Format("02 Jan 2006 15:04 WIB"))
	case "LessThan1Hr":
		body = fmt.Sprintf("Hello,\n\nYour task '%s' is due in less than 1 hour! Please complete is ASAP.\n\nDeadline: %s\n", todo.Title, todo.EndDate.Format("02 Jan 2006 15:04 WIB"))
	case "Overdue":
		body = fmt.Sprintf("Hello,\n\nYour task '%s' has passed its deadline and is still incomplete!\n\nOriginal Deadline: %s\n\nPlease complete it ASAP.\n", todo.Title, todo.EndDate.Format("02 Jan 2006 15:04 WIB"))
	default:
		body = fmt.Sprintf("Hello,\n\nThis is a reminder for your task:'%s'.\n", todo.Title)
	}

	msg := []byte(subject + "\n" + body)

	err := smtp.SendMail(host+":"+port, auth, from, []string{to}, msg)
	if err != nil {
		return fmt.Errorf("failed to send reminder email: %w", err)
	}

	return nil
}
