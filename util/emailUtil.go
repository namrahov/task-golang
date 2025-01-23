package util

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"task-golang/config"
	"task-golang/model"
)

var (
	mailChangeUrl     = "https://example.com/change-email"
	forgetPasswordUrl = "https://example.com/forget-password"
)

func GenerateActivationEmail(token string, linkType model.LinkType) model.EmailDto {
	email := model.EmailDto{
		From:    config.Props.UserFrom,
		Subject: "Activation email",
	}

	var link, body string

	switch linkType {
	case model.Registration:
		link = config.Props.UserActivationUrl + "?token=" + token
		body = fmt.Sprintf("Hormetli, istifadeci. Saytimizda qeydiyyatdan kecdiyiniz ucun tesekkur edirik. "+
			"Profilinizi aktivlesdirmek ucun bu linke kecin %s", link)
	case model.ChangeEmail:
		link = mailChangeUrl + "?token=" + token
		body = changeEmail(link)
	case model.ForgetPassword:
		link = forgetPasswordUrl + "?token=" + token
		body = forgetPasswordEmail(link)
	case model.SetPassword:
		body = setPasswordEmail(token)
	}

	email.Body = body

	fmt.Println("activation email =", email)
	return email
}

func changeEmail(link string) string {
	return fmt.Sprintf("To change your email, click the following link: %s", link)
}

func forgetPasswordEmail(link string) string {
	return fmt.Sprintf("To reset your password, click the following link: %s", link)
}

func setPasswordEmail(token string) string {
	return fmt.Sprintf("Use the following token to set your password: %s", token)
}

// SendEmailAsync sends an email asynchronously
func SendEmailAsync(from, to, subject, body string) {
	go func() {
		mailer := gomail.NewMessage()
		mailer.SetHeader("From", from)
		mailer.SetHeader("To", to)
		mailer.SetHeader("Subject", subject)
		mailer.SetBody("text/plain", body)

		dialer := gomail.NewDialer("smtp.gmail.com", 587, from, config.Props.UserPassword)
		dialer.SSL = false

		err := dialer.DialAndSend(mailer)
		if err != nil {
			fmt.Printf("Failed to send email: %v\n", err)
		} else {
			fmt.Println("Email sent successfully!")
		}
	}()
}
